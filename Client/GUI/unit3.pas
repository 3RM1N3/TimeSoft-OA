unit Unit3;

{$mode objfpc}{$H+}

interface

uses
  Classes, SysUtils, Forms, Controls, Graphics, Dialogs, Menus, ComCtrls,
  StdCtrls, ExtCtrls, ValEdit, Types;

type

  { TMainForm }

  TMainForm = class(TForm)
    ListView1: TListView;
    PageControl1: TPageControl;
    StatusBar1: TStatusBar;
    PageScan: TTabSheet;
    PageEditPic: TTabSheet;
    PageRework: TTabSheet;
    procedure FormCreate(Sender: TObject);
    procedure GroupBox1Click(Sender: TObject);
    procedure PageControl1Change(Sender: TObject);
    procedure PageScanContextPopup(Sender: TObject; MousePos: TPoint;
      var Handled: Boolean);
  private

  public

  end;

var
  MainForm: TMainForm;

implementation

{$R *.lfm}

{ TMainForm }

procedure TMainForm.GroupBox1Click(Sender: TObject);
begin

end;

procedure TMainForm.FormCreate(Sender: TObject);
begin

end;

procedure TMainForm.PageControl1Change(Sender: TObject);
begin

end;

procedure TMainForm.PageScanContextPopup(Sender: TObject; MousePos: TPoint;
  var Handled: Boolean);
begin

end;

end.

