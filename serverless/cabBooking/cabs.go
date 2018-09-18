package main

import "sort"

type CabResponse struct {
	Company  string  `json:"company"`
	Cab      string  `json:"cab"`
	Estimate float64 `json:"Estimate"`
	Arriving int     `json:"arriving"`
}

type CabList []CabResponse
type SortBy func(p1, p2 *CabResponse) bool

type cabsListSorter struct {
	cabs   CabList
	sortBy func(p1, p2 *CabResponse) bool // Closure used in the Less method.
}

func (sortBy SortBy) Sort(cabs CabList, orderBy string) {
	ps := &cabsListSorter{
		cabs:   cabs,
		sortBy: sortBy, // The sort method's receiver is the function (closure) that defines the sort order.
	}
	if orderBy == "highest" {
		sort.Sort(sort.Reverse(ps))
	} else {
		sort.Sort(ps)
	}
}

// Len is part of sort.Interface.
func (s *cabsListSorter) Len() int {
	return len(s.cabs)
}

// Swap is part of sort.Interface.
func (s *cabsListSorter) Swap(i, j int) {
	s.cabs[i], s.cabs[j] = s.cabs[j], s.cabs[i]
}

// Less is part of sort.Interface. It is implemented by calling the "sortBy" closure in the sorter.
func (s *cabsListSorter) Less(i, j int) bool {
	return s.sortBy(&s.cabs[i], &s.cabs[j])
}
