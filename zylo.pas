{*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22
 * License : GNU General Public License v3 (see LICENSE)
 * Author: Journal of Hamradio Informatics (http://pafelog.net)
*******************************************************************************}
unit UzLogExtension;

interface

uses
	Classes,
	Windows,
	Dialogs,
	IOUtils,
	UITypes,
	StrUtils,
	SysUtils,
	dprocess,
	UzLogQSO,
	UzLogGlobal,
	RegularExpressions;

type
	TzLogEvent = (evAddQSO = 0, evModifyQSO, evDeleteQSO);
	TImportDialog = class(TOpenDialog)
		procedure ImportMenuClicked(Sender: TObject);
	end;
	TExportDialog = class(TSaveDialog)
		procedure ExportMenuClicked(Sender: TObject);
		procedure FilterTypeChanged(Sender: TObject);
	end;

var
	Fmt: string;
	Handle: THandle;
	ImportDialog: TImportDialog;
	ExportDialog: TExportDialog;
	zlaunch: procedure(str: PAnsiChar); stdcall;
	zrevise: procedure(ptr: pointer; len: integer); stdcall;
	zverify: function (ptr: pointer; len: integer): integer; stdcall;
	zresult: function (ptr: pointer; len: integer): integer; stdcall;
	zinsert: procedure(ptr: pointer; len: integer); stdcall;
	zdelete: procedure(ptr: pointer; len: integer); stdcall;
	zfinish: procedure(); stdcall;

const
	BIN = 'zbin';
	LEN = sizeof(TQSOData);

(*zLog event handlers*)
procedure zLogInitialize();
procedure zLogContestInit(contest, cfg: string);
procedure zLogContestEvent(event: TzLogEvent; bQSO, aQSO: TQSO);
procedure zLogContestTerm();
procedure zLogTerminate();
function zLogCalcPointsHookHandler(aQSO: TQSO): boolean;
function zLogExtractMultiHookHandler(aQSO: TQSO; var mul: string): boolean;
function zLogValidMultiHookHandler(mul: string; var val: boolean): boolean;
function zLogGetTotalScore(): integer;

implementation

uses
	main;

procedure zLogInitialize();
var
	filter: AnsiString;
begin
	ImportDialog := TImportDialog.Create(MainForm);
	ExportDialog := TExportDialog.Create(MainForm);
	ImportDialog.Filter := '';
	ExportDialog.Filter := '';
	ExportDialog.OnTypeChange := ExportDialog.FilterTypeChanged;
	ExportDialog.Options := [ofOverwritePrompt];
	MainForm.MergeFile1.Caption := '&Import...';
	RunCommand('qxsl', ['filter'], filter, [poNoConsole]);
	if filter <> '' then begin
		ImportDialog.Filter := string(filter);
		ExportDialog.Filter := string(filter);
		MainForm.MergeFile1.OnClick := ImportDialog.ImportMenuClicked;
		MainForm.Export1.OnClick    := ExportDialog.ExportMenuClicked;
	end;
end;

procedure zLogContestInit(contest: string; cfg: string);
var
	dll: string;
	txt: string;
begin
	txt := TFile.ReadAllText(cfg);
	txt := TRegEx.Replace(txt, '(?m);.*?$', '');
	txt := TRegEx.Replace(txt, '(?m)#.*?$', '');
	dll := TRegEx.Replace(cfg, '(?i)\.CFG', '.DLL');
	if TRegEx.Match(txt, '(?im)^ *DLL +ON *$').Success then
	try
		Handle := LoadLibrary(PChar(dll));
		zlaunch := GetProcAddress(Handle, 'zlaunch');
		zrevise := GetProcAddress(Handle, 'zrevise');
		zverify := GetProcAddress(Handle, 'zverify');
		zresult := GetProcAddress(Handle, 'zresult');
		zinsert := GetProcAddress(Handle, 'zinsert');
		zdelete := GetProcAddress(Handle, 'zdelete');
		zfinish := GetProcAddress(Handle, 'zfinish');
		zlaunch(PAnsiChar(AnsiString(cfg)));
	except
		Handle := 0;
		MessageDlg(dll + ' not found.', mtWarning, [mbOK], 0);
	end;
end;

procedure zLogContestEvent(event: TzLogEvent; bQSO, aQSO: TQSO);
var
	qso: TQSOData;
begin
	if Handle <> 0 then begin
		if event <> evDeleteQSO then begin
			qso := aQSO.FileRecord;
			zinsert(@qso, LEN);
		end;
		if event <> evAddQSO then begin
			qso := bQSO.FileRecord;
			zdelete(@qso, LEN);
		end;
	end;
end;

procedure zLogContestTerm();
begin
	if Handle <> 0 then begin
		zfinish();
		FreeLibrary(Handle);
		Handle := 0;
	end;
end;

procedure zLogTerminate();
begin
end;

(*returns whether the QSO score is calculated by this handler*)
function zLogCalcPointsHookHandler(aQSO: TQSO): boolean;
var
	qso: TQSOData;
begin
	Result := Handle <> 0;
	if Result then begin
		qso := aQSO.FileRecord;
		aQSO.Points := zverify(@qso, LEN);
	end;
end;

(*returns whether the multiplier is extracted by this handler*)
function zLogExtractMultiHookHandler(aQSO: TQSO; var mul: string): boolean;
var
	qso: TQSOData;
begin
	Result := Handle <> 0;
	if Result then begin
		qso := aQSO.FileRecord;
		zrevise(@qso, LEN);
		mul := string(qso.Multi1);
	end;
end;

(*returns whether the multiplier is validated by this handler*)
function zLogValidMultiHookHandler(mul: string; var val: boolean): boolean;
begin
	Result := Handle <> 0;
	if Result then
		val := mul <> '';
end;

function zLogGetTotalScore(): Integer;
var
	buf: TBytesStream;
	qso: TQSOData;
	idx: integer;
begin
	Result := -1;
	if Handle <> 0 then begin
		buf := TBytesStream.Create;
		try
			for idx := 1 to Log.TotalQSO do begin
				qso := Log.QsoList[idx].FileRecord;
				buf.Write(qso, LEN);
			end;
			Result := zresult(buf.bytes, Log.TotalQSO);
		finally
  		buf.Free;
    end;
	end;
end;

procedure TImportDialog.ImportMenuClicked(Sender: TObject);
var
	tmp: string;
	msg: AnsiString;
begin
	if Execute then
	try
		tmp := TPath.GetTempFileName;
		RunCommand('qxsl', ['format', FileName, tmp, BIN], msg, [poNoConsole]);
		Log.QsoList.MergeFile(tmp, True);
		Log.SortByTime;
		MyContest.Renew;
		MainForm.EditScreen.RefreshScreen;
	finally
		TFile.Delete(tmp);
  end;
end;

procedure TExportDialog.ExportMenuClicked(Sender: TObject);
var
	tmp: string;
	msg: AnsiString;
begin
	FilterTypeChanged(Sender);
	FileName := ChangeFileExt(CurrentFileName, DefaultExt);
	if Execute then begin
		tmp := TPath.GetTempFileName;
		Log.SaveToFile(tmp);
		RunCommand('qxsl', ['format', tmp, FileName, Fmt], msg, [poNoConsole]);
  end;
end;

procedure TExportDialog.FilterTypeChanged(Sender: TObject);
var
	ext: string;
begin
	ext := SplitString(Filter, '|')[2 * FilterIndex - 1];
	Fmt := SplitString(Filter, '|')[2 * FilterIndex - 2];
	ext := TRegEx.Split(ext, ';')[0];
	ext := copy(ext, 2, Length(ext));
	DefaultExt := ext;
end;

initialization

finalization

end.
