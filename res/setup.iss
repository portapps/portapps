#define appName "@APP_NAME@™"
#define appHomepage "@APP_HOMEPAGE@"
#define pappId = "@PAPP_ID@"
#define pappGuid = "{@PAPP_GUID@"
#define pappName "@APP_NAME@™ Portable"
#define pappVersion "@PAPP_VERSION@"
#define pappURL "@PAPP_URL@"
#define pappFolder "@PAPP_FOLDER@"
#define pappPublisher "@PAPP_PUBLISHER@"
#define currentYear GetDateTimeString('yyyy', '', '');

[Setup]
AppId={#pappGuid}
AppName={#pappName}
AppVersion={#pappVersion}
;AppVerName={#pappName} {#pappVersion}
AppPublisher={#pappPublisher}
AppPublisherURL={#pappURL}
AppSupportURL={#pappURL}
AppUpdatesURL={#pappURL}

WizardImageFile=setup.bmp
WizardSmallImageFile=setup-mini.bmp
DisableWelcomePage=no
ShowLanguageDialog=yes
SetupIconFile=papp.ico

Compression=@ISS_COMPRESSION@
SolidCompression=yes

DefaultDirName={sd}\portapps\{#pappId}
CreateAppDir=yes
Uninstallable=no
PrivilegesRequired=lowest
ArchitecturesAllowed=@ISS_ARCH@

VersionInfoCompany={#pappPublisher}
VersionInfoCopyright={#pappPublisher} {#currentYear}
VersionInfoProductName={#pappName}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "@BIN_PATH@\build\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs

[InstallDelete]
Type: filesandordirs; Name: "{app}\{#pappFolder}"
Type: filesandordirs; Name: "{app}\log\*.log"
Type: filesandordirs; Name: "{app}\*.log"
Type: filesandordirs; Name: "{app}\{#pappId}*.exe"

#include "@BIN_PATH@\setup\run.iss"

[Messages]
WelcomeLabel1=[name] {#pappVersion}
WelcomeLabel2=[name] assists you with running {#appName} in a portable way without needing to install it in Windows.%n%nIf the application is running, it is recommended to close it before continuing.

[Code]

var
  IsUpgrade: Boolean;
  NoticePage: TOutputMsgMemoWizardPage;
  NoticeAcceptedRadio: TRadioButton;
  NoticeNotAcceptedRadio: TRadioButton;

procedure CheckNoticeAccepted(Sender: TObject);
begin
  WizardForm.NextButton.Enabled := NoticeAcceptedRadio.Checked;
end;

function CloneNoticeRadioButton(Source: TRadioButton): TRadioButton;
begin
  Result := TRadioButton.Create(WizardForm);
  Result.Parent := NoticePage.Surface;
  Result.Left := Source.Left;
  Result.Top := Source.Top;
  Result.Width := Source.Width;
  Result.Height := Source.Height;
  Result.OnClick := @CheckNoticeAccepted;
end;

procedure InitializeWizard();
begin
  NoticePage := CreateOutputMsgMemoPage(wpWelcome,
    'Notice of Non-Affiliation and Disclaimer', 'Please read the following important information before continuing.',
    'When you are ready to continue with Setup, click Next.',
    '');

  NoticePage.RichEditViewer.Height := WizardForm.LicenseMemo.Height;
  NoticePage.RichEditViewer.Text :=
    'Portapps is not affiliated, associated, authorized, endorsed by, or in any way officially connected with {#appName}, or any of its subsidiaries or its affiliates.' + #13#10 + #13#10 +
    'The official {#appName} website can be found at {#appHomepage}.' + #13#10 + #13#10 +
    'The name {#appName} as well as related names, marks, emblems and images are registered trademarks of their respective owners.';

  NoticeAcceptedRadio := CloneNoticeRadioButton(WizardForm.LicenseAcceptedRadio);
  NoticeAcceptedRadio.Caption := 'I understand this notice';
  NoticeNotAcceptedRadio := CloneNoticeRadioButton(WizardForm.LicenseNotAcceptedRadio);
  NoticeNotAcceptedRadio.Caption := 'I''d rather not continue';

  NoticeNotAcceptedRadio.Checked := True;
end;

procedure CurPageChanged(CurPageID: Integer);
begin
  if CurPageID = NoticePage.ID then
  begin
    CheckNoticeAccepted(nil);
  end;
end;

function isEmptyDir(dirName: String): Boolean;
var
  FindRec: TFindRec;
  FileCount: Integer;
begin
  Result := False;
  if FindFirst(dirName + '\*', FindRec) then begin
    try
      repeat
        if (FindRec.Name <> '.') and (FindRec.Name <> '..') then begin
          FileCount := 1;
          break;
        end;
      until not FindNext(FindRec);
    finally
      FindClose(FindRec);
      if FileCount = 0 then Result := True;
    end;
  end;
end;

function InitializeSetup: Boolean;
begin
  Result := True;
  IsUpgrade := False;
end;

function NextButtonClick(PageId: Integer): Boolean;
begin
  Result := True;
  if (PageId = wpSelectDir) then begin
    if DirExists(ExpandConstant('{app}')) and not isEmptyDir(ExpandConstant('{app}')) then begin
      IsUpgrade := True;
    end;
    if IsUpgrade then begin
      if FileExists(ExpandConstant('{app}\portapp.json')) then begin
        Result := True;
      end
      else begin
        MsgBox(ExpandConstant('The selected dir is not empty or you are performing an upgrade of {#pappName} but required file portapp.json cannot be found. Please select a correct folder.'), mbError, MB_OK);
        Result := False;
        exit;
      end;
    end;
  end;
end;
