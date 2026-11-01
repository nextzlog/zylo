# ZyLO

ZyLO is the plugin DLL interface of [zLog](https://use.zlog.org) 2.8.3.0 or later.
zLog loads each plugin DLL at start-up, calls the functions it exports and gives it callbacks into zLog.
A plugin can score a user-defined contest, add import and export formats, and drive the zLog window.
The reference plugin is [zlog-max.dll](zlog-max.md).

## Loading

zLog reads the DLL names from the `DLLs` key of the `zylo` section in `zlog.ini` and loads them from the plugin directory.

```ini
[zylo]
DLLs=zlog-max.dll
```

The plugin manager writes this key when the user installs a DLL package.
zLog calls `zylo_finish_event` on exit and unloads the DLL with `FreeLibrary` if the function returns true.

## Conventions

The DLL is a 64-bit Windows DLL.
Every function uses the `stdcall` convention, which is the standard x64 calling convention.

| Type | Meaning |
|---|---|
| `int` | 32-bit signed integer, `integer` in Delphi. |
| `bool` | 1-byte boolean, 0 for false. |
| `ptr` | pointer to a QSO record. |
| `in` | NUL-terminated string from zLog, in the ANSI code page of Windows (Shift_JIS in Japan). |
| `out` | NUL-terminated string to zLog, in UTF-8. |
| `fn` | pointer to a callback function. |

A string or a record passed to the DLL is valid only during the call.
The DLL must not unwind a panic or an exception into zLog.

## Sequence

```text
zLog                                      DLL
 | zylo_record_size(384)    (optional)  -> |
 | zylo_allow_*(callback)   (optional)  -> |
 | zylo_query_format(callback)          -> |  registers formats
 | zylo_launch_event()                  -> |  true to stay loaded
 |                                         |
 | contest opened                          |
 | zylo_offset_event(minutes)           -> |
 | zylo_attach_event(contest, cfg)      -> |
 | zylo_assign_event(contest, cfg)      -> |  only the scoring DLL
 | zylo_query_cities(callback)          -> |  only the scoring DLL
 | zylo_insert_event / delete_event     -> |  every log change
 | zylo_verify_event(record)            -> |  only the scoring DLL
 | zylo_points_event(points, multis)    -> |  only the scoring DLL
 | zylo_detach_event(contest, cfg)      -> |
 |                                         |
 | zylo_finish_event()                  -> |  true to be unloaded
```

## Required Functions

zLog refuses a DLL that lacks one of these functions.

| Function | Signature | Called |
|---|---|---|
| `zylo_launch_event` | `bool ()` | Once after loading. Return false to report a version mismatch. |
| `zylo_finish_event` | `bool ()` | Once on exit. Return true to let zLog unload the DLL. |
| `zylo_window_event` | `void (MSG *msg)` | For each window message of zLog. |
| `zylo_offset_event` | `void (int minutes)` | Before `zylo_attach_event`, with the UTC offset of the log in minutes. |
| `zylo_attach_event` | `void (in contest, in cfg)` | When a contest is opened, with its name and the path to its CFG file. |
| `zylo_assign_event` | `void (in contest, in cfg)` | After `zylo_attach_event`, only for the DLL named by the CFG file. |
| `zylo_detach_event` | `void (in contest, in cfg)` | When the contest is closed. |
| `zylo_insert_event` | `void (ptr record)` | After a QSO is added, and for the new record of an edited QSO. |
| `zylo_delete_event` | `void (ptr record)` | After a QSO is deleted, and for the old record of an edited QSO, before the insert. |
| `zylo_verify_event` | `void (ptr record)` | Before a QSO is recorded, several times, only for the scoring DLL. |
| `zylo_points_event` | `int (int points, int multis)` | When the score is shown, only for the scoring DLL. |
| `zylo_query_cities` | `void (fn cities)` | When the multiplier list is built, only for the scoring DLL. |
| `zylo_import_event` | `bool (in source, in target)` | When the user imports a file in a registered format. |
| `zylo_export_event` | `bool (in target, in format)` | When the user exports a file in a registered format. |
| `zylo_button_event` | `void (int handler, int button)` | When a button or a menu item with a registered handler is clicked. |
| `zylo_editor_event` | `void (int handler, int key)` | When a key is typed into an editor with a registered handler. |

zLog sends insert and delete events to every DLL, not only to the scoring DLL.
A DLL that keeps its own copy of the log reads the existing QSOs from the file returned by the query `{F}`.

## Optional Functions

zLog calls each of these functions only if the DLL exports it.

| Function | Signature | Called |
|---|---|---|
| `zylo_record_size` | `int (int size)` | Once after loading, with 384. Return 384 to receive ZLOX records, any other value for ZLO records. |
| `zylo_query_format` | `void (fn format)` | Once after loading, to register file formats. |
| `zylo_allow_insert` | `void (fn insert)` | Once after loading. |
| `zylo_allow_delete` | `void (fn delete)` | Once after loading. |
| `zylo_allow_update` | `void (fn update)` | Once after loading. |
| `zylo_allow_commit` | `void (fn commit)` | Once after loading. |
| `zylo_allow_dialog` | `void (fn dialog)` | Once after loading. |
| `zylo_allow_notify` | `void (fn notify)` | Once after loading. |
| `zylo_allow_access` | `void (fn access)` | Once after loading. |
| `zylo_allow_handle` | `void (fn handle)` | Once after loading. |
| `zylo_allow_button` | `void (fn button)` | Once after loading. |
| `zylo_allow_editor` | `void (fn editor)` | Once after loading. |
| `zylo_allow_script` | `void (fn script)` | Once after loading. |

## Callbacks

The DLL keeps the pointers given by the optional functions and calls them while zLog runs.
Records passed to the insert, delete, update and commit callbacks are ZLO records of 256 bytes.

| Callback | Signature | Effect |
|---|---|---|
| `insert` | `void (ptr record)` | Adds the QSO to the log. |
| `delete` | `void (ptr record)` | Deletes the QSO from the log. |
| `update` | `void (ptr record)` | Replaces the QSO with the same ID. |
| `commit` | `void (ptr record, int id, uint64 mask)` | Copies the fields chosen by `mask` from the record to the QSO with the given ID, without sending events back to the DLL. |
| `dialog` | `void (out text)` | Shows the text in a dialog. |
| `notify` | `void (out text)` | Shows the text in a toast. |
| `access` | `void (char *query)` | Evaluates the query in a buffer of 257 bytes and overwrites it with the answer of at most 256 bytes. |
| `handle` | `HANDLE (out name)` | Returns the window handle of a zLog component, or 0. |
| `button` | `int (out name)` | Registers a handler for a button or a menu item and returns its number. |
| `editor` | `int (out name)` | Registers a handler for an editor and returns its number. |
| `script` | `int (out expression)` | Runs a Delphi expression and returns its ordinal value, 1 for other values, or 0 on error. |
| `format` | `void (out filter)` | Adds file formats to the import and export dialogs. |
| `cities` | `void (out list)` | Fills the multiplier list. |

The bits of the commit mask are as follows.

| Bit | Field |
|---|---|
| 0 | call sign |
| 1 | sent number |
| 2 | received number |
| 3 | operator |
| 4 | memo |
| 5 | multiplier 1 |
| 6 | multiplier 2 |

The filter has the form of a Windows file dialog filter, such as `ADIF (*.adi)|*.adi|Cabrillo (*.cbr)|*.cbr`.
Only the last DLL that registers a filter handles import and export.
The cities list has one multiplier per line, the code and then the name, separated by spaces.

## Queries

The `access` callback evaluates the variables of the CW keyboard of zLog, such as `$B`, `$X`, `$R`, `$F`, `$Z`, `$I`, `$Q`, `$V`, `$O`, `$S`, `$P`, `$A`, `$N`, `$L`, `$C`, `$E` and `$M`.
It also evaluates the following variables.

| Query | Answer |
|---|---|
| `{V}` | version of zLog |
| `{F}` | path to the ZLO or ZLOX file of the log |
| `{C}` | own call sign |
| `{B}` | operating band |
| `{M}` | operating mode |

## Scoring

The CFG file of a user-defined contest names its scoring DLL by a `dll` line before `exit`.

```text
dll zlog-max.dll
exit
script lua
```

zLog Max ignores the line for `zlog-max.dll` and scores the same CFG file by itself, with or without the Lua part.

zLog asks the scoring DLL instead of its own rules for the points, the multipliers and the total of the contest.
If the named DLL is not installed, zLog opens the plugin manager at its package.

`zylo_verify_event` receives the QSO before it is recorded and writes the result into the same record.
The DLL sets the points and the multipliers and leaves the other fields as they are.
An empty multiplier 1 marks the QSO as invalid.
`zylo_points_event` receives the sums computed by zLog and returns the total score.

## Import and Export

zLog imports a file by calling `zylo_import_event` with the chosen file and a temporary file.
The DLL converts the chosen file into a ZLO or ZLOX file at the temporary path and returns true, and zLog merges it into the log.
zLog exports a file by first saving the log in its own format at the chosen path.
It then calls `zylo_export_event` with that path and the name of the chosen format, and the DLL converts the file in place.
The name is the part of the filter before `|`, such as `ADIF (*.adi)`.

## Records

A record is one QSO of a ZLO file, 256 bytes in little endian.
A DLL that returned 384 from `zylo_record_size` receives a ZLOX record, which appends 128 bytes to the same 256 bytes.
A short string is a length byte followed by its characters in the ANSI code page; the size below includes the length byte.
The QSO ID is `id / 100` of the field at offset 88.

| Offset | Size | Field |
|---|---|---|
| 0 | 8 | time, `TDateTime` (days since 1899-12-30) |
| 8 | 13 | call sign, short string |
| 21 | 31 | sent number, short string |
| 52 | 31 | received number, short string |
| 83 | 1 | reserved |
| 84 | 2 | sent RST |
| 86 | 2 | received RST |
| 88 | 4 | serial number and ID |
| 92 | 1 | mode |
| 93 | 1 | band |
| 94 | 1 | power |
| 95 | 31 | multiplier 1, short string |
| 126 | 31 | multiplier 2, short string |
| 157 | 1 | new multiplier 1 |
| 158 | 1 | new multiplier 2 |
| 159 | 1 | points |
| 160 | 15 | operator, short string |
| 175 | 65 | memo, short string |
| 240 | 1 | CQ |
| 241 | 1 | dupe |
| 242 | 1 | reserved |
| 243 | 1 | transmitter number |
| 244 | 4 | power 2 |
| 248 | 4 | reserved 2 |
| 252 | 4 | reserved 3 |

The ZLOX extension follows.

| Offset | Size | Field |
|---|---|---|
| 256 | 11 | frequency, short string |
| 267 | 1 | QSY violation |
| 268 | 11 | PC name, short string |
| 279 | 1 | forced |
| 280 | 1 | QSL state |
| 281 | 1 | invalid |
| 282 | 1 | verified by RBN |
| 283 | 3 | continent, short string |
| 286 | 11 | entity, short string |
| 297 | 31 | multiplier 3, short string |
| 328 | 1 | new multiplier 3 |
| 329 | 55 | reserved |

The modes are numbered as follows.

| Number | Mode |
|---|---|
| 0 | CW |
| 1 | SSB |
| 2 | FM |
| 3 | AM |
| 4 | RTTY |
| 5 | other |

The bands are numbered as follows.

| Number | Band |
|---|---|
| 0 | 1.9 MHz |
| 1 | 3.5 MHz |
| 2 | 7 MHz |
| 3 | 10 MHz |
| 4 | 14 MHz |
| 5 | 18 MHz |
| 6 | 21 MHz |
| 7 | 24 MHz |
| 8 | 28 MHz |
| 9 | 50 MHz |
| 10 | 144 MHz |
| 11 | 430 MHz |
| 12 | 1200 MHz |
| 13 | 2400 MHz |
| 14 | 5600 MHz |
| 15 | 10 GHz and up |
