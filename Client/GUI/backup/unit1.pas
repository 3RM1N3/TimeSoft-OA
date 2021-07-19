unit Unit1;

{$mode objfpc}{$H+}

interface

uses
  Classes, SysUtils, Forms, Controls, Graphics, Dialogs, StdCtrls;

type

  { TFormLogin }

  TFormLogin = class(TForm)
    ButtonLogin: TButton;
    Edit1: TEdit;
    EditUser: TEdit;
    EditServerAddr: TEdit;
    Label1: TLabel;
    Label2: TLabel;
    Label3: TLabel;
    Label4: TLabel;
    LabelSignUp: TLabel;
    procedure FormCreate(Sender: TObject);
    procedure Label1Click(Sender: TObject);
  private

  public

  end;

var
  FormLogin: TFormLogin;

implementation

{$R *.lfm}

{ TFormLogin }

procedure TFormLogin.Label1Click(Sender: TObject);
begin

end;

procedure TFormLogin.FormCreate(Sender: TObject);
begin

end;

end.

