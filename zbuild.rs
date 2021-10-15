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

type Return<E> = Result<E, Box<dyn Error>>;

fn init(pkg: &str) -> Return<String> {
	Cmd::new("go").args(["mod", "init", pkg]).status()?;
	Cmd::new("go").args(["get", "-u", "all"]).status()?;
	Cmd::new("go").args(["mod", "tidy"]).status()?;
	Ok(format!("{}.dll", pkg))
}

fn make(pkg: &str) -> Return<()> {
	let mode = "-buildmode=c-shared";
	let args = ["build", "-o", &init(pkg)?, mode];
	let make = Cmd::new("go").args(&args).output()?;
	io::stdout().write_all(&make.stdout)?;
	io::stderr().write_all(&make.stderr)?;
	std::process::exit(make.status.code().unwrap());
}

fn name(path: &path::Path) -> Option<String> {
	Some(path.file_name()?.to_str()?.to_string())
}

fn save(name: &str, data: &[u8]) {
	let mut opts = fs::OpenOptions::new();
	match opts.create_new(true).write(true).open(&name) {
		Ok(mut file) => file.write_all(data).unwrap_or(()),
		Err(_err) => eprintln!("{} already exists.", name),
	}
}

fn item(item: &mut Value) -> Return<String> {
	let mal = "malformed TOML file";
	let val = item.as_table_mut().ok_or(mal)?;
	let url = val["url"].as_str().ok_or(mal)?;
	let bin = get(&url.to_string())?.bytes()?;
	let sum = format!("{:x}", md5::compute(bin));
	val.insert("sum".to_string(), Value::String(sum));
	Ok(format!("{}\n", item))
}

fn kind(list: &mut Value) -> Return<String> {
	let mal = "malformed TOML file";
	let items = list.as_table_mut();
	for (_, it) in items.ok_or(mal)? {
		item(it)?;
	}
	Ok(format!("{}\n", list))
}

fn table(mut list: Value) -> Return<String> {
	let mal = "malformed TOML file";
	let items = list.as_table_mut();
	for (_, it) in items.ok_or(mal)? {
		kind(it)?;
	}
	Ok(format!("{}\n", list))
}

fn fetch(url: &str) -> Return<String> {
	let spec = include_str!("schema.yaml");
	let temp = serde_yaml::from_str(spec)?;
	let test = JSONSchema::compile(&temp).unwrap();
	let toml = get(url)?.text()?.parse::<Value>()?;
	let json = serde_json::to_value(toml.clone())?;
	if let Err(error) = test.validate(&json) {
		eprintln!("{}", join(error, ", "));
		exit(1);
	}
	table(toml)
}

fn merge() -> Return<String> {
	let head = include_str!("master.toml");
	let list = include_str!("market.list");
	let mut toml = head.to_string();
	for url in list.lines() {
		toml.push_str(&fetch(url)?);
	}
	Ok(toml)
}

#[argopt::subcmd]
pub fn markets() -> Return<()> {
	let source = merge()?;
	let target = io::stdout();
	let mut de = Deserializer::new(&source);
	let mut en = Serializer::pretty(target);
	Ok(transcode(&mut de, &mut en)?)
}

#[argopt::subcmd]
fn compile() -> Return<()> {
	let error = "failed to determine project name";
	save("zutils.go", include_bytes!("zutils.go"));
	make(&name(&env::current_dir()?).ok_or(error)?)
}

#[argopt::cmd_group(commands = [compile, markets])]
fn main() -> Return<()> {
	env::set_var("GOOS", "windows");
	env::set_var("GOARCH", "amd64");
	env::set_var("CGO_ENABLED", "1");
	env::set_var("CC", "x86_64-w64-mingw32-gcc");
}
