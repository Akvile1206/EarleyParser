package main

import (
	"fmt"
)

var rules = map[string][][]string{
	"S": {{"NP", "VP"}},
	"NP": {
		{"N"},
		{"N", "PP"},
	},
	"PP": {{"P","NP"}},
	"VP": {
		{"VP","PP"},
		{"V","VP"},
		{"V","NP"},
		{"V"},
	},
}

var terminals = map[string][]string{
	"N" : {"can", "fish", "they", "rivers", "December"},
	"P" : {"in"},
	"V" : {"can", "fish"},
}

type progress struct {
	rule   string;
	before []string;
	after  []string;
}

type entry struct {
	rule  progress;
	start int;
	end   int;
	hist  []int;
	id    int;
	phase string;
}

func (e *entry) notFinished() bool {
	return (len(e.rule.after) != 0);
}

func (e *entry) predictable() bool {
	if !e.notFinished() {
		return false;
	}
	N := e.rule.after[0];
	_, ok := terminals[N];
	if ok {
		return false;
	}
	return true;
}

func (e *entry) addTo(step int) {
	if e.duplicate(step) {
		return;
	}
	e.id = id;
	id++;
	c.steps[step][c.entries[step]] = *e;
	c.entries[step]++;
}

func (e *entry) duplicate(step int) bool{
	for i := 0; i < c.entries[step]; i++ {
		f := c.steps[step][i];
		a := (e.rule.rule == f.rule.rule);
		b := equal(e.rule.before, f.rule.before);
		c := equal(e.rule.after, f.rule.after);
		d := (e.start == f.start) && (e.end == f.end);
		e := equalInt(e.hist, f.hist);
		if a && b && c && d && e {
			return true;
		}
	}
	return false;
}

func equal(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func equalInt(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

type chart struct {
	steps [][]entry;
	entries []int;
}

var c = chart {};
var id = 0;

func print(steps int) {
	fmt.Println("STEP|ID | RULE          |(start,end)|PHASE| HIST     ");
	for i := 0; i < steps; i++ {
		for j := 0; j < c.entries[i]; j++ {
			printEntry(c.steps[i][j], i);
		} 
	}
}

func printEntry(e entry, i int) {
	b := "";
	for _, s := range(e.rule.before) {
		b += s + " ";
	}
	a := "";
	for _, s := range(e.rule.after) {
		a += " "+s;
	}
	s := b + "." + a;
	fmt.Printf("%4d|%3d|%2v->%11v| (%4d,%2d) |  %1v  | %v\n",
		i,
		e.id,
		e.rule.rule, 
		s, 
		e.start, 
		e.end,
		e.phase,
		e.hist,
	)
}

func (e *entry) predict(step int) {
	N := e.rule.after[0];
	rulz := rules[N]; 
	for _, body := range(rulz) {
		e := entry{
			progress{
				N,
				[]string{},
				body,
			},
			e.end,
			e.end,
			[]int{},
			0,
			"P",
		};
		e.addTo(step);
	}
}

func (e *entry) scan(word string, step int) {
	N := e.rule.after[0];
	words := terminals[N];
	for _, w := range(words) {
		if w == word {
			e := entry {
				progress {N, []string{word}, []string{}},
				e.start,
				e.end + 1,
				[]int{},
				0,
				"S",
			};
			e.addTo(step + 1);
		}
	}
}

func (e *entry) complete(step int) {
	for _, s := range(c.steps) {
		for _, f := range(s) {
			if len(f.rule.after) == 0 {
				continue;
			}
			if f.rule.after[0] != e.rule.rule {
				continue;
			} 
			if f.end != e.start {
				continue;
			}
			new_entry := entry {
				progress{
					f.rule.rule,
					append(f.rule.before, e.rule.rule),
					f.rule.after[1:],
				},
				f.start,
				e.end,
				append(f.hist, e.id),
				0,
				"C",
			}
			new_entry.addTo(step);
		}
	} 
}

func main() {
	words := []string{"they", "can", "fish", "in", "rivers","in", "December", "$"};
	c.steps = make([][]entry, len(words) + 1);
	c.entries = make([]int, len(words) + 1);
	for i := 0; i < len(words) + 1; i++ {
		c.entries[i] = 0;
		c.steps[i] = make([]entry, 50);
	}
	//initialise
	e := entry{
		progress{"S",[]string{},[]string{"NP","VP"}},
		0,
		0,
		[]int{},
		0,
		"C",
	}
	e.addTo(0);
	for i, w := range(words) {
		for j := 0; j < c.entries[i]; j++ {
			e := c.steps[i][j];
			if (e.notFinished()) {
				if (e.predictable()) {
					e.predict(i);
				} else {
					e.scan(w, i)
				}
			} else {
				e.complete(i);
			}
		}
	}
	print(len(words));
}