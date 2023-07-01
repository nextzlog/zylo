/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22nd
 * Released under the MIT License (or GPL v3 until 2021 Oct 28th) (see LICENSE)
 * Univ. Tokyo Amateur Radio Club Development Task Force (https://nextzlog.dev)
*******************************************************************************/

const OPT: &str = "-replace zylo=github.com/nextzlog/zylo/src/commons@HEAD";

use minijinja::render;
use std::io::Write;
use std::process::abort;
use std::process::Command as Cmd;
use std::{env, fs, path};

type Return<E> = Result<E, Box<dyn std::error::Error>>;

fn ok(code: i32) {
	if code != 0 {
		abort();
	}
}

fn init(pkg: &str) -> Return<String> {
	shell("go", &format!("mod init {}", pkg));
	shell("go", &format!("mod edit {}", OPT));
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
	fs::write(md5, format!("{v:x}", v = sum))?;
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

#[allow(unused_must_use)]
fn shell(cmd: &str, arg: &str) {
	let seq = arg.split_whitespace();
	Cmd::new(cmd).args(seq).status();
}

#[argopt::subcmd]
fn compile(#[opt(default_value = "2.8")] ver: String) -> Return<()> {
	let lib = render!(include_str!("zbuild.go"), version => ver);
	save("main.go", lib.as_bytes());
	make(&name(&env::current_dir()?).unwrap())
}

#[argopt::subcmd]
fn setup() -> Return<()> {
	shell("apt", "install -y gcc-mingw-w64 golang-go upx");
	shell("brew", "install mingw-w64 go upx");
	shell("choco", "install mingw golang upx");
	shell("pacman", "-Sy mingw-w64-gcc go upx");
	Ok(())
}

#[argopt::cmd_group(commands = [compile, setup])]
fn main() -> Return<()> {
	env::set_var("GOOS", "windows");
	env::set_var("GOARCH", "amd64");
	env::set_var("CGO_ENABLED", "1");
	env::set_var("GOPROXY", "direct");
	env::set_var("CC", "x86_64-w64-mingw32-gcc");
}
