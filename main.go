package main

import (
	"fmt"

	//"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	mb "mBoxGo/mbox"
)


const hma string = "<span foreground=\"black\" background=\"white\" size=\"medium\"><tt><b>"
const hmb string = "</b></tt></span>"
const dma string = "<span foreground=\"blue\" background=\"white\" size=\"medium\"><tt>"
const dmb string = "</tt></span>"
const cma string = "<span foreground=\"black\" background=\"yellow\" size=\"medium\"><tt>"
const cmb string = "</tt></span>"

func main() {

	mbData := [][]string{
		{"Date", "Name", "URL", "Info"},
		{"20190904", "Vodafone", "www.vodafone.com", "Mi cuenta en la web de Vodafone"},
		{"20191001", "micuenta", "gmail.com", "Cuenta de correo en gmail"},
		{"20190522", "BNK", "www.banco.com", "Pues eso, el banco y tal"},
		{"20181105", "García", "www.zaragoza.es", "Ejemplo con tilde, y alguna cosilla"}}

	gtk.Init(nil)

		fmt.Println(mbData)
		
		mbox := mb.NewMBox(mbData, true, []string{" "})
		
		fmt.Println(mbox)

	gtk.Main()

}


/*
void main(string[] args) {

	string[][] mbData = [
		["Date", "Name", "URL", "Info"],
		["20190904", "Vodafone", "www.vodafone.com", "Mi cuenta en la web de Vodafone"],
		["20191001", "micuenta", "gmail.com", "Cuenta de correo en gmail"],
		["20190522", "BNK", "www.banco.com", "Pues eso, el banco y tal"],
		["20181105", "García", "www.zaragoza.es", "Ejemplo con tilde, y alguna cosilla"]];

	MainWin mwin;

	Main.init(args);

	mwin = new MainWin(mbData);

	Main.run();

}

class MainWin : MainWindow {
	Keymap keymap;
	MBox mbox;
	auto rnd = Random(42);

	this(string[][] mbData) {
		super("mBox Test");
		addOnDestroy(delegate void(Widget w) { mainQuit(); });

		setDefaultSize(600, 400);
		keymap = Keymap.getDefault();
		addOnKeyPress(&onKeyPress);

		auto aligns = ["rigth", "left", "center", "left"];
		mbox = new MBox(mbData, true, aligns);
		mbox.setCursorMarkup(cma, cmb);

		addOnScroll(delegate bool(Event e, Widget w) {
			writeln("scroll event");
			return true;
		});

		mbox.addOnButtonPress(delegate bool(Event e, Widget w) {
			auto eb = e.button();

			if (e.isDoubleClick(eb)) {
				writeln("mbox double check: get row data");
				if (mbox.activeData() != []) {
					writeln(mbox.activeData());
				} else {
					writeln("no data active");
				}

			} else {
				//writeln("mbox single check: get position");
			}
			return true;
		});

		add(new MainBox(mbData, mbox, this));
		showAll();
	}

	void mainQuit() {
		Main.quit();
		writeln("Bye.");
	}

	bool onKeyPress(GdkEventKey* eventKey, Widget widget) {
		string key = keymap.keyvalName(eventKey.keyval);

		switch (key) {
			case "Up":
				mbox.cursorUp();
				break;
			case "Down":
				mbox.cursorDown();
				break;
			case "Escape":
				if (mbox.cursorIsActive()) {
					mbox.clearCursor();
				} else {
					mainQuit();
				}
				break;
			case "Return":
				if (mbox.activeData() != []) {
					writeln(mbox.activeData());
				} else {
					writeln("no data active");
				}
				break;
			case "Delete":
				mbox.deleteActiveRow();
				break;
			case "Insert":
				mbox.addRow(modify(["20190101", "Mi veloz router",
					"www.here.com", "Acceso all router de casa"]));
				break;
			case "F12":
				mbox.reverseData();
				break;
			case "e":
				if (eventKey.state & ModifierType.CONTROL_MASK) {
					auto toEdit = mbox.activeData();
					if ( toEdit != []) {
						auto edited = modify(toEdit);
						mbox.editActiveRow(edited);
					}
				}
				break;
			default:
				writeln("New: ", key);
				break;
		}

		return true;
	}

	private string[] modify(string[] str) {
		auto r = uniform(0, 255, rnd);
		auto y = uniform(2000, 2020, rnd);
		auto m = uniform(1, 13, rnd);
		auto d = uniform(1, 29, rnd);
		auto date = to!string(y) ~ "/" ~ to!string(m) ~ "/" ~ to!string(d);
		auto url = "www." ~ to!string(r) ~ ".com";
		str[0] = date;
		str[2] = url;
		return str;
	}
}

class MainBox : Box {
	this(string[][] mbData, MBox mbox, MainWin mwin) {
		super(Orientation.VERTICAL, 0);

		auto headText = new Label("Ejemplo");
		headText.setMarkup("<span foreground=\"green\"><b>Ejemplo</b></span>");
		headText.setMarginTop(8);
		add(headText);

		add(mbox);

		auto close = new Button("Close");
		packEnd(close, false, false, 0);

		close.addOnClicked(delegate void(Button b) { mwin.mainQuit(); });
	}
}
*/
