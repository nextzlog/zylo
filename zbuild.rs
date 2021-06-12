/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
 */

use std::env::current_dir;
use std::env::set_var;
use std::env::var;
use std::env::vars_os;
use std::fs::OpenOptions;
use std::io::Error;
use std::io::Write;
use std::io::stderr;
use std::io::stdout;
use std::path::Path;
use std::process::Command;
use std::process::exit;

fn setenv() -> String {
	set_var("GOOS", "windows");
	set_var("GOARCH", "amd64");
	set_var("CGO_ENABLED", "1");
	set_var("CC", "x86_64-w64-mingw32-gcc");
	for (key, value) in vars_os() {
		println!("{:?}: {:?}", key, value);
	}
	return var("REPO").unwrap_or("zylo/dll".to_string());
}

fn output(name: String, data: &[u8]) {
	match OpenOptions::new()
		.create_new(true)
		.write(true)
		.open(&name) {
		Ok(mut file) => file.write_all(data).unwrap_or(()),
		Err(e) => eprintln!("{} not saved by {}", name, e),
	}
}

fn main() -> Result<(), Error> {
	let repo = setenv();
	let path = current_dir()?;
	let file = Path::new(&path).file_name();
	let name = file.unwrap().to_str().unwrap();
	output("qxsl.exe".to_string(), include_bytes!("qxsl.exe"));
	output("qxsl.fmt".to_string(), include_bytes!("qxsl.fmt"));
	output("asset.go".to_string(), include_bytes!("assets/asset.go"));
	output(format!("{}.go", name), include_bytes!("assets/toast.go"));
	Command::new("go").arg("mod").arg("init").arg(repo).status()?;
	Command::new("go").arg("get").arg("-u").arg("all").status()?;
	Command::new("go").arg("mod").arg("tidy").status()?;
	let arg = format!("build -o {d}.dll -buildmode=c-shared", d = name);
	let out = Command::new("go").args(arg.split_whitespace()).output()?;
	stdout().write_all(&out.stdout).unwrap();
	stderr().write_all(&out.stderr).unwrap();
	exit(out.status.code().unwrap_or(1));
}
