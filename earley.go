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
	step  int;
}

var chart = make([]entry, 1000);
var id = 0;

func print() {
	fmt.Println("STEP| ID | RULE          |(start,end)| HIST     ");
	for i := 0; i < id; i++ {
		printEntry(chart[i], i);
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
	fmt.Printf("%4d|%3d|%2v->%11v| (%4d,%2d) | %v\n",
		e.step,
		i, 
		e.rule.rule, 
		s, 
		e.start, 
		e.end,
		e.hist,
	)
}

func predict(step int) {
	predicted := make(map[string]bool);
	for i := 0; i < id; i++ {
		if chart[i].step != step {
			continue;
		} 
		if len(chart[i].rule.after) == 0 {
			continue;
		}
		N := chart[i].rule.after[0];
		_, ok := predicted[N];
		if ok {
			continue;
		}
		predicted[N] = true;
		_, ok = terminals[N];
		if ok {
			continue;
		}
		//not yet predicted and not a special terminal node
		rulz := rules[N]; 
		for _, body := range(rulz) {
			chart[id] = entry{
				progress{
					N,
					[]string{},
					body,
				},
				step,
				step,
				[]int{},
				step,
			};
			id++;
		}
	}
}

func scan(word string, step int) {
	scaned := make(map[string]bool);
	for i := 0; i < id; i++ {
		if chart[i].step != step {
			continue;
		} 
		if len(chart[i].rule.after) == 0 {
			continue;
		}
		N := chart[i].rule.after[0];
		_, ok := scaned[N];
		if ok {
			continue;
		}
		scaned[N] = true;
		var words []string;
		words, ok = terminals[N];
		if !ok {
			continue;
		}
		for _, w := range(words) {
			if w == word {
				chart[id] = entry {
					progress {N, []string{word}, []string{}},
					chart[i].start,
					step + 1,
					[]int{},
					step,
				};
				id++;
			}
		}
	}
}

func complete(j int, step int) {
	found := true;
	completed := make(map[int]bool);
	for found {
		found = false;
		for i := 0; i < id; i++ {
			if chart[i].step != j && chart[i].step != step {
				continue;
			}
			if(len(chart[i].rule.after) > 0) {
				continue;
			}
			_, ok := completed[i];
			if ok {
				continue;
			}
			completed[i] = true;
			found = true;
			rule := chart[i].rule;
			for a := 0; a < id; a++ {
				if chart[a].end != chart[i].start {
					continue;
				}
				if len(chart[a].rule.after) == 0 {
					continue;
				}
				if chart[a].rule.after[0] != rule.rule {
					continue;
				}
				var e = entry{
					progress{
						chart[a].rule.rule,
						append(chart[a].rule.before, rule.rule),
						chart[a].rule.after[1:],
					},
					chart[a].start,
					chart[i].end,
					append(chart[a].hist, i),
					step,
				}
				if (!duplicate(e)) {
					chart[id] = e;
					id++;
				}
				
			}
		}
	} 
}

func duplicate(e entry) bool{
	for i := 0; i < id; i++ {
		f := chart[i];
		a := (e.rule.rule == f.rule.rule);
		b := (len(e.rule.before) == len(f.rule.before));
		c := (len(e.rule.after) == len(f.rule.after));
		d := (e.start == f.start) && (e.end == f.end);
		e := equal(e.hist, f.hist);
		if a && b && c && d && e {
			return true;
		}
	}
	return false;
}

func equal(a, b []int) bool {
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

func main() {
	//initialise
	chart[id] = entry{
		progress{"S",[]string{},[]string{"NP","VP"}},
		0,
		0,
		[]int{},
		0,
	}
	id++;
	predict(0);
	scan("they", 0);
	complete(0, 1);
	predict(1);
	scan("can", 1);
	complete(1, 2);
	predict(2);
	scan("fish", 2);
	complete(2, 3);
	predict(3);
	scan("in", 3);
	complete(3, 4);
	predict(4);
	scan("rivers", 4);
	complete(4, 5);
	predict(5);
	scan("in", 5);
	complete(5, 6);
	predict(6);
	scan("December", 6);
	complete(6, 7);
	print();
}