/***************************************************************************
* Copyright (c) 2013 Nikita Mikhailov <nslqqq@gmail.com>
*
* Permission is hereby granted, free of charge, to any person
* obtaining a copy of this software and associated documentation
* files (the "Software"), to deal in the Software without restriction,
* including without limitation the rights to use, copy, modify, merge,
* publish, distribute, sublicense, and/or sell copies of the Software,
* and to permit persons to whom the Software is furnished to do so,
* subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included
* in all copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
* THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
* OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
* ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE
* OR OTHER DEALINGS IN THE SOFTWARE.
*
***************************************************************************/

import QtQuick 2.0

QtObject {
    readonly property string capslockWarning:   qsTr("Warnung: Feststelltaste ist aktiv")
    readonly property string layout:            qsTr("Layout")
    readonly property string login:             qsTr("Login")
    readonly property string loginFailed:       qsTr("Login fehlgeschlagen")
    readonly property string loginSucceeded:    qsTr("Login erfolgreich")
    readonly property string password:          qsTr("Passwort eingeben")
    readonly property string emptyPassword:     qsTr("Bitte ein Passwort eingeben!")
    readonly property string passwordChange:    qsTr("Passwort ändern")
    readonly property string prompt:            qsTr("Geben Sie Ihren Benutzernamen und Ihr Passwort ein")
    readonly property string promptSelectUser:  qsTr("Wählen Sie Ihren Benutzer und geben Sie Ihr Passwort ein")
    readonly property string promptUser:        qsTr("Benutzernamen eingeben")
    readonly property string promptPassword:    qsTr("Passwort eingeben")
    readonly property string emptyPrompt:       qsTr("Passwort:")
    readonly property string showPasswordPrompt:qsTr("Passwort anzeigen")
    readonly property string hidePasswordPrompt:qsTr("Passwort verstecken")
    readonly property string reboot:            qsTr("Neustart")
    readonly property string session:           qsTr("Sitzung")
    readonly property string shutdown:          qsTr("Herunterfahren")
    readonly property string suspend:           qsTr("Schlafen")
    readonly property string hibernate:         qsTr("Ruhen")
    readonly property string userName:          qsTr("Benutzername")
    readonly property string welcomeText:       qsTr("Willkommen zu %1")
    readonly property string pamMaxtriesError:  qsTr("Passwortänderung wurde abgebrochen, weil die maximale Anzahl an Versuchen erreicht wurde")
    readonly property string pamMaxtriesInfo:   qsTr("Neues Passwort ist rund! Bitte geben Sie Ihr aktuelles Passwort erneut ein!")
}

