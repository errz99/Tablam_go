package main

/*
import std.stdio;
import std.typecons;
import std.conv;
import std.algorithm;
import std.uni : byGrapheme;
import std.array : replicate;
import std.range.primitives : walkLength;

import gtk.Widget;
import gtk.Grid;
import gtk.Box;
import gtk.Label;
import gtk.EventBox;
import gtk.Button;
import glib.ListG;
import gdk.Event;

const string hma = "<span><tt><b>";
const string hmb = "</b></tt></span>";
const string dma = "<span background=\"white\"><tt>";
const string dmb = "</tt></span>";
const string cma = "<span foreground=\"white\" background=\"#6666dd\"><tt>";
const string cmb = "</tt></span>";

class MBox : Grid {
private:
	bool _hasHead;
	RowBox[] rbs;
	string[][] data;
	string[][] datax;
	int position;
	int outPosition = -1;
	int lastPosition = -1;
	string[2] headMarkup;
	string[2] dataMarkup;
	string[2] cursorMarkup;
	int hsep = 3;
	ulong[] max;
	int[] changedMax;
	int separation = 1;
	auto sep = " ";
	string[] aligns;

public:
	this(string[][] data, bool hasHead, string[] als = []) {
		headMarkup = [hma, hmb];
		dataMarkup = [dma, dmb];
		cursorMarkup = [cma, cmb];

		_hasHead = hasHead;
		if (_hasHead == true) outPosition++;

		position = outPosition;
		lastPosition = outPosition;
		sep = " ".replicate(separation);
		max.length = data[0].length;

		if (als == []) {
			aligns.length = data[0].length;
			foreach (ref ha; aligns) {
				ha = "right";
			}
		} else {
			aligns = als;
		}

		super();
		setHalign(Align.CENTER);
		setBorderWidth(8);
		setRowSpacing(hsep);

		foreach (d; data) {
			addRow(d);
		}
	}

	void setHeadMarkup(string a, string b) {
		headMarkup = [a, b];
	}

	void setDataMarkup(string a, string b) {
		dataMarkup = [a, b];
	}

	void setCursorMarkup(string a, string b) {
		cursorMarkup = [a, b];
	}

	void setElemAlign(int i, string halign) {
		aligns[i] = halign;
	}

	void cursorDown() {
		if ((_hasHead && rbs.length > 1) || (!_hasHead && rbs.length > 0)) {
			lastPosition = position;

			if (position < rbs.length - 1) {
				position++;
			} else {
				position = outPosition + 1;
			}

			updateCursor();
		}
		writeln("down ", position);
	}

	void cursorUp() {
		lastPosition = position;
		position--;

		if (position < outPosition + 1) position = to!int(rbs.length) - 1;
		if (position >= 0) updateCursor();
	}

	bool cursorIsActive() {
		if (position > outPosition) {
			return true;
		} else {
			return false;
		}
	}

	void clearCursor() {
		if (position > outPosition) {
			for (int i = 0; i < rbs[0].labels.length; i++) {
				rbs[position].labels[i].setMarkup(
					dataMarkup[0] ~ rbs[position].datax[i] ~ dataMarkup[1]);
			}
			position = outPosition;
		}
	}

	string[] activeData() {
		if (position > outPosition) {
			return rbs[position].data;
		} else {
			return [];
		}
	}

	void editActiveRow(string[] edata) {
		writeln(edata);
		changedMax = [];
		auto edatax = newX(edata);

		rbs[position].data = edata;
		rbs[position].datax = edatax;
		data[position] = edata;
		datax[position] = edatax;

		updateChanged();
		markupActiveRow();
	}

	void deleteActiveRow() {
		if (position > outPosition && position < rbs.length) {

			rbs = rbs[0..position] ~ rbs[position + 1..$];
			data = data[0..position] ~ data[position + 1..$];
			datax = datax[0..position] ~ datax[position + 1..$];

			removeRow(position);

			for (int i = 0; i < rbs.length; i++) {
				rbs[i].setName(to!string(i));
			}

			if (rbs.length == outPosition + 1) {
				position = outPosition;
			} else if (position == rbs.length) {
				position--;
			}

			if (position > outPosition) {
				markupActiveRow();
			}
		}
	}

	void reverseData() {
		if (_hasHead) {
			reverse(data[1..$]);
			reverse(datax[1..$]);

		} else {
			reverse(data);
			reverse(datax);
		}

		for (int i = outPosition + 1; i < datax.length; i++) {
			rbs[i].data = data[i];
			rbs[i].datax = datax[i];
			for (int j = 0; j < rbs[i].labels.length; j++) {
				applyMarkup(i, j, rbs[i].datax[j]);
			}
		}
	}

	void addRow(string[] rdata) {
		auto rb = new RowBox(rdata, this);
		changedMax = rb.changedMax;
		rbs ~= rb;
		attach(rb, 0, cast(int)datax.length, 1, 1);
		data ~= rb.data;
		datax ~= rb.datax;
		updateChanged();
	}

	private void updateCursor() {
		if (lastPosition > outPosition) {
			for (int i = 0; i < rbs[0].labels.length; i++) {
				rbs[lastPosition].labels[i].setMarkup(
					dataMarkup[0] ~ rbs[lastPosition].datax[i] ~ dataMarkup[1]);
			}
		}
		if (position > outPosition) {
			for (int i = 0; i < rbs[0].labels.length; i++) {
				rbs[position].labels[i].setMarkup(
					cursorMarkup[0] ~ rbs[position].datax[i] ~ cursorMarkup[1]);
			}
		}
	}

	private void updateChanged() {
		foreach (cm; changedMax) {
			int j;
			while (j < rbs.length) {
				auto elemgr = rbs[j].data[cm].byGrapheme;
				ulong grow = max[cm] - elemgr.walkLength;
				auto elemx = createX(rbs[j].data[cm], cm, grow);

				rbs[j].datax[cm] = elemx;
				applyMarkup(j, cm, elemx);
				++j;
			}
		}
	}

	private void markupActiveRow() {
		for (int i = 0; i < rbs[0].labels.length; i++) {
			rbs[position].labels[i].setMarkup(
				cursorMarkup[0] ~ rbs[position].datax[i] ~ cursorMarkup[1]);
		}
	}

	private string createX(string elem, int i, ulong grow) {
		if (aligns[i] == "left") {
			return sep ~ elem ~ " ".replicate(grow) ~ sep;

		} else if (aligns[i] == "rigth") {
			return sep ~ " ".replicate(grow) ~ elem ~ sep;

		} else if (aligns[i] == "center") {
			ulong a = grow / 2;
			ulong b = grow / 2;
			if (grow % 2 != 0) { b++; }
			return sep ~ " ".replicate(a) ~ elem ~ " ".replicate(b) ~ sep;

		} else {
			return sep ~ elem ~ " ".replicate(grow) ~ sep;
		}
	}

	private void applyMarkup(int i, int j, ref string elemx) {
		if (_hasHead == true && i == 0) {
			rbs[i].labels[j].setMarkup(headMarkup[0] ~ elemx ~ headMarkup[1]);
		} else if (i == position) {
			rbs[i].labels[j].setMarkup(cursorMarkup[0] ~ elemx ~ cursorMarkup[1]);
		} else {
			rbs[i].labels[j].setMarkup(dataMarkup[0] ~ elemx ~ dataMarkup[1]);
		}
	}

	private string[] newX(string[] ndata) {
		string[] ndatax;

		for (int i = 0; i < ndata.length; i++) {
			auto elemgr = ndata[i].byGrapheme;

			if (elemgr.walkLength > max[i]) {
				max[i] = elemgr.walkLength;
				changedMax ~= i;
			}

			ulong grow = max[i] - elemgr.walkLength;
			ndatax ~= createX(ndata[i], i, grow);
		}

		return ndatax;
	}
}

class RowBox : Box {
	string[] data;
	string[] datax;
	Label[] labels;
	int[] changedMax;

	this(string[] d, MBox mb) {
		auto idx = cast(int) mb.rbs.length;
		super(Orientation.HORIZONTAL, mb.hsep);
		setName(to!string(idx));
		data = d.dup;

		//datax = mb.newX(data);

		for (int i = 0; i < data.length; i++) {
			auto elemgr = data[i].byGrapheme;

			if (elemgr.walkLength > mb.max[i]) {
				mb.max[i] = elemgr.walkLength;
				changedMax ~= i;
			}

			ulong grow = mb.max[i] - elemgr.walkLength;
			datax ~= mb.createX(data[i], i, grow);
		}

		foreach (ref elemx; datax) {
			auto ebox = new EventBox();
			add(ebox);
			auto label = new Label(elemx);
			label.setMarkup(mb.dataMarkup[0] ~ elemx ~ mb.dataMarkup[1]);
			ebox.add(label);
			labels ~= label;
		}

		addOnButtonPress(delegate bool(Event e, Widget w) {
			auto eb = e.button();
			auto name = getName();

			if (to!int(name) > mb.outPosition) {
				if (e.isDoubleClick(eb)) {
					//
				} else if (mb.position != to!int(name)) {
					mb.lastPosition = mb.position;
					mb.position = to!int(name);
					mb.updateCursor();
				}
			} else {
				//writeln("button pressed at header");
			}
			return false;
		});

		showAll();
	}
}
*/
