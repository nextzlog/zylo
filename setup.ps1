[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12
wget http://github.com/nextzlog/qxsl/releases/download/nightly/qxsl.exe -OutFile qxsl.exe
git clone https://github.com/jr8ppg/zLog -b v2510RT
git clone https://github.com/z505/TProcess-Delphi
$dir0="zLog/zlog"
$dir1="zLog/zlog/Win64/Debug"
$dir2="zLog/zlog/Win64/Release"
New-Item $dir1 -ItemType Directory
New-Item $dir2 -ItemType Directory
Copy-Item TProcess-Delphi/dpipes.pas $dir0
Copy-Item TProcess-Delphi/dprocess.pas $dir0
Copy-Item TProcess-Delphi/pipes_win.inc $dir0
Copy-Item TProcess-Delphi/process_win.inc $dir0
Copy-Item qxsl.exe $dir1
Copy-Item qxsl.exe $dir2
Copy-Item zylo.pas zLog/zlog/UzLogExtension.pas
Invoke-Item zLog/zlog/Zlog.dpr
