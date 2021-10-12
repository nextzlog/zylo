/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
 */

use itertools::join;
use jsonschema::JSONSchema;
use reqwest::blocking::get;
use serde_json::Serializer;
use serde_transcode::transcode;
use std::error::Error;
use std::io::Write;
use std::process::exit;
use std::process::Command as Cmd;
use std::{env, fs, io, path};
use toml::Deserializer;
use toml::Value;

type CommandResult<E> = Result<E, Box<dyn Error>>;

const CSHARED: &str = "-buildmode=c-shared";
const MALFORM: &str = "malformed TOML file";

fn save(dir: &path::Path, name: &str, data: &[u8]) {
	let mut opts = fs::OpenOptions::new();
	opts.create_new(true).write(true);
	match opts.open(dir.join(&name)) {
		Ok(mut file) => file.write_all(data).unwrap_or(()),
		Err(e) => eprintln!("{} is not saved {}", name, e),
	}
}

fn check(mut table: Value) -> CommandResult<String> {
	for (_class, items) in table.as_table_mut().ok_or(MALFORM)? {
		for (_name, item) in items.as_table_mut().ok_or(MALFORM)? {
			let val = item.as_table_mut().ok_or(MALFORM)?;
			let url = val["url"].as_str().ok_or(MALFORM)?;
			let bin = get(&url.to_string())?.bytes()?;
			let sum = format!("{:x}", md5::compute(bin));
			val.insert("sum".to_string(), Value::String(sum));
		}
	}
	Ok(table.to_string())
}

fn fetch(url: &str) -> CommandResult<String> {
	let spec = include_str!("schema.yaml");
	let temp = serde_yaml::from_str(spec)?;
	let test = JSONSchema::compile(&temp).unwrap();
	let toml = get(url)?.text()?.parse::<Value>()?;
	let json = serde_json::to_value(toml.clone())?;
	if let Err(error) = test.validate(&json) {
		eprintln!("{}", join(error, ", "));
		exit(1);
	}
	check(toml)
}

fn merge() -> CommandResult<String> {
	let mut toml = String::new();
	for url in include_str!("market.list").lines() {
		toml.push_str(&format!("{}\n", fetch(url)?));
	}
	Ok(toml)
}

#[argopt::subcmd]
fn compile() -> CommandResult<()> {
	let dir = env::current_dir()?.canonicalize()?;
	let lib = dir.file_name().unwrap().to_str().unwrap();
	save(&dir, "zutils.go", include_bytes!("zutils.go"));
	Cmd::new("go").arg("mod").arg("init").arg(lib).status()?;
	Cmd::new("go").arg("get").arg("-u").arg("all").status()?;
	Cmd::new("go").arg("mod").arg("tidy").status()?;
	let dll = &format!("{}.dll", lib);
	let arg = ["build", "-a", "-o", dll, CSHARED];
	let cmd = Cmd::new("go").args(&arg).output()?;
	io::stdout().write_all(&cmd.stdout)?;
	io::stderr().write_all(&cmd.stderr)?;
	std::process::exit(cmd.status.code().unwrap());
}

#[argopt::subcmd]
fn markets() -> CommandResult<()> {
	let source = merge()?;
	let target = io::stdout();
	let mut de = Deserializer::new(&source);
	let mut en = Serializer::pretty(target);
	Ok(transcode(&mut de, &mut en)?)
}

#[argopt::cmd_group(commands = [compile, markets])]
fn main() -> CommandResult<()> {
	env::set_var("GOOS", "windows");
	env::set_var("GOARCH", "amd64");
	env::set_var("CGO_ENABLED", "1");
	env::set_var("CC", "x86_64-w64-mingw32-gcc");
}
