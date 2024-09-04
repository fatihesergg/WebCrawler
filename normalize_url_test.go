package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
        // add more test cases here
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestUrlsFromHTML(t *testing.T) {
    tests := []struct { 
        inputURL      string
        inputBody     string 
        expectedURLS  []string
        expectedError error
    }{
        {
            inputURL: "https://blog.boot.dev",
            inputBody: `
            <html>
            <body>
            <a href="/path/one">
            <span>Boot.dev</span>
            </a>
            /body>
            </html>
            `,
            expectedURLS: []string{
                "https://blog.boot.dev/path/one",
            },
            expectedError: nil,

        },

        {

            inputURL: "https://boot.dev",
            inputBody: `
            <html>
            <body>
            <a href="https://boot.dev/path/two">
            <span>Boot.dev</span>
            </a>
            </body>
            </html>
            `,
            expectedURLS: []string{
                "https://boot.dev/path/two",
            },
            expectedError: nil,
        },
        {

            inputURL: "https://boot.dev",
            inputBody: `<html><body><a href="/u/fatihesergg" class="w-full hover:opacity-75"><div data-v-d717fdb3="" class="min-h-12 min-w-12"><div data-v-d717fdb3="" class="square-container"><img data-v-d717fdb3="" loading="lazy" src="https://orbitermag.com/wp-content/uploads/2017/03/default-user-image-300x300.png" alt="user avatar" class="absolute left-1/4 top-1/4 w-1/2 scale-125 rounded-full object-cover" style="transform: scale(1.75);"><img data-v-d717fdb3="" loading="lazy" src="https://www.boot.dev/_nuxt/4.xMhV6BS_.png" alt="user avatar" class="absolute object-cover" style="transform: scale(1.75);"></div></div></a></body></html>`,
            expectedURLS: []string{
                "https://boot.dev/u/fatihesergg",
            },
            expectedError: nil,
        },
        {
            inputURL: "https://boot.dev",
            inputBody: `<html><body><a href="/path/three"></a></body></html>`,
            expectedURLS: []string{
                "https://boot.dev/path/three"},
            expectedError: nil,
        },
        
    }
    for i_,tc := range tests {
        t.Run(tc.inputURL, func(t *testing.T) {
         actualURLS, err := getURLsFromHTML(tc.inputBody, tc.inputURL) 
        if err != nil {
            t.Errorf("Test %v - %s FAIL: unexpected error: %v", i_, tc.inputURL, err) 
        }
        if len(actualURLS) != len(tc.expectedURLS) {
            t.Errorf("Test %v - %s FAIL: length of expected URLS: %v, length of actual: %v",i_, tc.inputURL,len(tc.expectedURLS), len(actualURLS)) 
        }

        for i, url := range actualURLS {
            if url != actualURLS[i]{
                t.Errorf("Test %v - %s FAIL: expected URLS: %v, actual: %v", i, tc.inputURL, tc.expectedURLS, actualURLS)
            }
        } 
   
        })
    }
}
