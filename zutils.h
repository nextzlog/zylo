/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
*/

typedef void (*InsertCallBack)(void*);
typedef void (*DeleteCallBack)(void*);
typedef void (*UpdateCallBack)(void*);
typedef void (*DialogCallBack)(char*);
typedef void (*AccessCallBack)(void*);
typedef long (*ButtonCallBack)(char*);
typedef long (*EditorCallBack)(char*);
typedef void (*FormatCallBack)(char*);
typedef void (*CitiesCallBack)(char*);

InsertCallBack insertCallBack;
DeleteCallBack deleteCallBack;
UpdateCallBack updateCallBack;
DialogCallBack dialogCallBack;
AccessCallBack accessCallBack;
ButtonCallBack buttonCallBack;
EditorCallBack editorCallBack;

inline void callInsert(void *qso) {insertCallBack(qso);}
inline void callDelete(void *qso) {deleteCallBack(qso);}
inline void callUpdate(void *qso) {updateCallBack(qso);}
inline void callDialog(char *str) {dialogCallBack(str);}
inline void callAccess(void *str) {accessCallBack(str);}
inline long callButton(char *str) {return buttonCallBack(str);}
inline long callEditor(char *str) {return editorCallBack(str);}
inline void callFormat(char *str, FormatCallBack cb) {cb(str);}
inline void callCities(char *str, CitiesCallBack cb) {cb(str);}
