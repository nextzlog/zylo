/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22
 * License: The MIT License since 2021 October 28 (see LICENSE)
 * Author: Journal of Hamradio Informatics (http://pafelog.net)
*******************************************************************************/

use itertools::join;
use jsonschema::JSONSchema;
use minijinja::{context, Environment, State};
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
use version_compare::Version;

type Return<E> = Result<E, Box<dyn Error>>;

fn ok(code: i32) {
	if code != 0 {
		abort();
	}
}

fn init(pkg: &str) -> Return<String> {
	shell("go", &format!("mod init {}", pkg));
	shell("go", &include_str!("go.mod.opts"));
	shell("go", "get -u all");
	shell("go", "mod tidy");
	Ok(format!("{}.dll", pkg))
}

fn make(pkg: &str) -> Return<()> {
	let name = &init(pkg)?;
	let args = ["build", "-o", &name, "-buildmode=c-shared"];
	ok(Cmd::new("go").args(&args).status()?.code().unwrap());
	shell("upx", name);
	let md5 = format!("{}.md5", name);
	let sum = md5::compute(fs::read(&name)?);
	fs::write(md5, format!("{v:x}", v=sum))?;
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

fn checksum(item: &mut Value) -> Return<()> {
	if item.get("sum").is_none() {
		let val = item.as_table_mut().unwrap();
		let url = val["url"].as_str().unwrap();
		let bin = get(url)?.error_for_status()?.bytes()?;
		let sum = format!("{:x}", md5::compute(bin));
		val.insert("sum".into(), Value::String(sum));
	}
	Ok(())
}

fn document(item: &mut Value) -> Return<()> {
	if item.get("doc").is_some() {
		let val = item.as_table_mut().unwrap();
		let url = val["doc"].as_str().unwrap();
		let txt = get(url)?.error_for_status()?.text()?;
		val.insert("doc".into(), Value::String(txt));
	}
	Ok(())
}

fn tree(list: &mut Value) -> Return<String> {
	let items = list.as_table_mut();
	for (_, it) in items.unwrap() {
		if it.is_table() {
			tree(it)?;
		}
		if it.get("url").is_some() {
			checksum(it)?;
		} else {
			document(it)?;
		}
	}
	Ok(list.to_string())
}

fn fetch(url: &str) -> Return<String> {
	let res = get(url)?.error_for_status()?;
	let val = res.text()?.parse::<Value>()?;
	let sch = from_str(include_str!("schema.yaml"))?;
	let cmp = JSONSchema::compile(&sch).unwrap();
	let tmp = serde_json::to_value(val.clone())?;
	if let Err(error) = cmp.validate(&tmp) {
		eprintln!("{}", join(error, ", "));
		ok(1);
	}
	tree(&mut val.clone())
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

fn older(_st: &State, now: String, old: String) -> bool {
	Version::from(&now).unwrap() < Version::from(&old).unwrap()
}

#[argopt::subcmd]
fn compile(#[opt(default_value = "2.8")]ver: String) -> Return<()> {
	let mut env = Environment::new();
	env.add_test("older_than", older);
	let src = include_str!("zutils.go");
	let ctx = context!(version => ver);
	let lib = env.render_str(src, ctx);
	save("zutils.go", lib?.as_bytes());
	make(&name(&env::current_dir()?).unwrap())
}

#[argopt::subcmd]
fn market() -> Return<()> {
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

#[argopt::cmd_group(commands = [compile, market, setup])]
fn main() -> Return<()> {
	env::set_var("GOOS", "windows");
	env::set_var("GOARCH", "amd64");
	env::set_var("CGO_ENABLED", "1");
	env::set_var("CC", "x86_64-w64-mingw32-gcc");
}
