package main

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"Struct with one string field",
			struct {
				Name string
			}{"Chris"},
			[]string{"Chris"},
		},
		{
			"Struct with two string fields",
			struct {
				Name string
				City string
			}{"Chris", "London"},
			[]string{"Chris", "London"},
		},
		{
			"Struct with non string fields",
			struct {
				Name string
				City string
				Age  int
			}{"Chris", "London", 33},
			[]string{"Chris", "London"},
		},
		{
			"Nested struct",
			Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"Pointer to things",
			&Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"slices",
			[]Profile{
				{33, "London"},
				{34, "Dubai"},
			},
			[]string{"London", "Dubai"},
		},
		{
			"Arrays",
			[2]Profile{
				{33, "London"},
				{34, "Dubai"},
			},
			[]string{"London", "Dubai"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})
	}

	t.Run("maps", func(t *testing.T) {
		aMap := map[string]string{
			"Foo": "Bar",
			"Baz": "Boz",
		}
		expected := []string{"Bar", "Boz"}

		var got []string
		walk(aMap, func(input string) {
			got = append(got, input)
		})

		for _, value := range expected {
			assertContains(t, got, value)
		}
	})

	t.Run("chan", func(t *testing.T) {
		aChannel := make(chan Profile)

		go func() {
			aChannel <- Profile{33, "Berlin"}
			aChannel <- Profile{34, "Paris"}
			close(aChannel)
		}()

		var got []string
		want := []string{"Berlin", "Paris"}

		walk(aChannel, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("function", func(t *testing.T) {
		aFunction := func() (Profile, Profile) {
			return Profile{33, "Berlin"}, Profile{34, "Paris"}
		}

		var got []string
		want := []string{"Berlin", "Paris"}

		walk(aFunction, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func assertContains(t *testing.T, haystack []string, needle string) {
	t.Helper()
	contains := false
	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}

	if !contains {
		t.Errorf("expected %+v to contain %q but it didn't", haystack, needle)
	}
}
