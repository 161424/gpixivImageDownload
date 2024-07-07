package model

type Ranks struct {
	User         string
	Day          []string
	Week         []string
	Month        []string
	Tops         int
	DownloadPath string
}

func (r Ranks) Content() string {
	s := ""
	if len(r.Day) != 0 {
		s += "Date=["
		for _, i := range r.Day {
			s += i
			s += ","
		}
	}
	s = s[:len(s)-1] + "];"
	if len(r.Week) != 0 {
		s += "Week=["
		for _, i := range r.Week {
			s += i
			s += ","
		}
	}
	s = s[:len(s)-1] + "];"
	if len(r.Month) != 0 {
		s += "Month=["
		for _, i := range r.Month {
			s += i
			s += ","
		}
	}
	s = s[:len(s)-1] + "];"
	return s
}
