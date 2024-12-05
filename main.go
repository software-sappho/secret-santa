package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

type Person struct {
	Name  string
	Gifts []string
}

type PageData struct {
	AssignedTo     string
	GiftSuggestion string
}

func main() {
	people := []Person{
		{Name: "Anna", Gifts: []string{"a break", "a one-week trip to Las Vegas", "some Casino tickets"}},
		{Name: "Bob", Gifts: []string{"a Blahaj from IKEA", "TRT", "weed"}},
		{Name: "Charlie", Gifts: []string{"a one-month subscription to a therapist", "an Adult Coloring Book", "Antidepressants"}},
		{Name: "David", Gifts: []string{"The Communist Manifesto", "a 4-day work week", "a Molotov Cocktail(ACAB)"}},
	}

	// Generate gift assignments (no self-gifting)
	assignments := generateAssignments(people)

	// Serve static files (like CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handle the main page and form submission
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodPost {
			// Handle form submission
			userName := r.FormValue("userName")
			assignedTo := ""
			for giver, receiver := range assignments {
				if giver == userName {
					assignedTo = receiver
					break
				}
			}

			// Suggest a gift for the recipient
			var suggestion string
			if assignedTo != "" {
				recipient := getRecipient(assignedTo, people)
				suggestion = getGiftSuggestions(recipient)
			}

			// Prepare data for the template
			data := PageData{
				AssignedTo:     assignedTo,
				GiftSuggestion: suggestion,
			}

			// Execute template with data
			tmpl.Execute(w, data)
		} else {
			// Just show the form if it's a GET request
			tmpl.Execute(w, nil)
		}
	})

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func generateAssignments(people []Person) map[string]string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	shuffled := make([]Person, len(people))
	copy(shuffled, people)
	random.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	assignments := make(map[string]string)

	for i := 0; i < len(people); i++ {
		if shuffled[i].Name == people[i].Name {
			return generateAssignments(people)
		}
		assignments[people[i].Name] = shuffled[i].Name
	}

	return assignments
}

func getRecipient(name string, people []Person) Person {
	for _, person := range people {
		if person.Name == name {
			return person
		}
	}
	return Person{}
}

func getGiftSuggestions(person Person) string {
	if len(person.Gifts) > 0 {
		return person.Gifts[rand.Intn(len(person.Gifts))]
	}
	return "No suggestions available"
}
