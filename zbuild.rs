/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
 */

use itertools::join;
use jsonschema::JSONSchema;
use reqwest::blocking::get;
use serde_json::Serializer;
use serde_transcode::transcode;
use serde_yaml::from_str;
use std::error::Error;
use std::io::Write;
use std::process::abort;
use std::process::Command as Cmd;
use std::{env, fs, io, path};
use toml::Deserializer;
use toml::Value;

type Return<E> = Result<E, Box<dyn Error>>;

fn ok(code: i32) {
	if code != 0 {
		abort();
	}
}

fn init(pkg: &str) -> Return<String> {
	Cmd::new("go").args(["mod", "init", pkg]).status()?;
	Cmd::new("go").args(["get", "-u", "all"]).status()?;
	Cmd::new("go").args(["mod", "tidy"]).status()?;
	Ok(format!("{}.dll", pkg))
}

fn make(pkg: &str) -> Return<()> {
	let name = &init(pkg)?;
	let args = ["build", "-o", &name, "-buildmode=c-shared"];
	ok(Cmd::new("go").args(&args).status()?.code().unwrap());
	shell("upx", name);
	Ok(())
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

fn leaf(item: &mut Value) -> Return<String> {
	let mal = "malformed TOML file";
	let val = item.as_table_mut().ok_or(mal)?;
	let url = val["url"].as_str().ok_or(mal)?;
	let bin = get(&url.to_string())?.bytes()?;
	let sum = format!("{v:x}", v = md5::compute(bin));
	val.insert("sum".to_string(), Value::String(sum));
	Ok(item.to_string())
}

fn tree(list: &mut Value) -> Return<String> {
	let mal = "malformed TOML file";
	let items = list.as_table_mut();
	for (_, it) in items.ok_or(mal)? {
		if it.get("url").is_some() {
			leaf(it)?;
		} else if it.is_table() {
			tree(it)?;
		}
	}
	Ok(list.to_string())
}

fn fetch(url: &str) -> Return<String> {
	let spec = from_str(include_str!("schema.yaml"))?;
	let test = JSONSchema::compile(&spec).unwrap();
	let toml = get(url)?.text()?.parse::<Value>()?;
	let json = serde_json::to_value(toml.clone())?;
	if let Err(error) = test.validate(&json) {
		eprintln!("{}", join(error, ", "));
		ok(1);
	}
	tree(&mut toml.clone())
}

fn merge() -> Return<String> {
	let mut toml = String::new();
	for url in include_str!("market.list").lines() {
		toml.push_str(&format!("{}\n", fetch(url)?));
	}
	Ok(toml)
}

#[allow(unused_must_use)]
fn shell(cmd: &str, arg: &str) {
	let seq = arg.split_whitespace();
	Cmd::new(cmd).args(seq).status();
}

#[argopt::subcmd]
fn compile() -> Return<()> {
	let error = "failed to determine project name";
	save("zutils.go", include_bytes!("zutils.go"));
	make(&name(&env::current_dir()?).ok_or(error)?)
}

#[argopt::subcmd]
fn markets() -> Return<()> {
	let source = merge()?;
	let target = io::stdout();
	let mut de = Deserializer::new(&source);
	let mut en = Serializer::pretty(target);
	Ok(transcode(&mut de, &mut en)?)
}

#[argopt::subcmd]
fn setup() -> Return<()> {
	let yaml = include_str!("launch.yaml");
	let data = from_str::<Value>(yaml)?;
	let cmds = data.as_table().unwrap();
	for (cmd, args) in cmds {
		shell(cmd, args.as_str().unwrap());
	}
	Ok(())
}

#[argopt::cmd_group(commands = [compile, markets, setup])]
fn main() -> Return<()> {
	env::set_var("GOOS", "windows");
	env::set_var("GOARCH", "amd64");
	env::set_var("CGO_ENABLED", "1");
	env::set_var("CC", "x86_64-w64-mingw32-gcc");
}
