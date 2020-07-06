package clinical

import (
	"gitlab.slade360emr.com/go/base"
	"reflect"
	"testing"
)

func TestHumanDateTime(t *testing.T) {
	type args struct {
		dt base.DateTime
	}
	tests := []struct {
		name string
		args args
		want base.Markdown
	}{
		{
			name: "YYYY",
			args: args{
				dt: base.DateTime("2020"),
			},
			want: base.Markdown("2020"),
		},
		{
			name: "YYYY-MM",
			args: args{
				dt: base.DateTime("2020-01"),
			},
			want: base.Markdown("Jan 2020"),
		},
		{
			name: "YYYY-MM-DD",
			args: args{
				dt: base.DateTime("2020-01-01"),
			},
			want: base.Markdown("Jan 01 2020"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HumanDateTime(tt.args.dt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HumanDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanInstant(t *testing.T) {
	type args struct {
		i base.Instant
	}
	tests := []struct {
		name string
		args args
		want base.Markdown
	}{
		{
			name: "good case",
			args: args{
				i: base.Instant("2015-02-07T13:28:17.239+02:00"),
			},
			want: base.Markdown("Sat, 07 Feb 2015 13:28:17 +0200"),
		},
		{
			name: "bad case, already formatted",
			args: args{
				i: base.Instant("Sat, 07 Feb 2015 13:28:17 +0200"),
			},
			want: base.Markdown("Sat, 07 Feb 2015 13:28:17 +0200"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HumanInstant(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HumanInstant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanNarrative(t *testing.T) {
	inputHTML := `
	<html>
	  <head>
		<title>My Mega Service</title>
		<link rel=\"stylesheet\" href=\"main.css\">
		<style type=\"text/css\">body { color: #fff; }</style>
	  </head>
	
	  <body>
		<div class="logo">
		  <a href="http://jaytaylor.com/"><img src="/logo-image.jpg" alt="Mega Service"/></a>
		</div>
	
		<h1>Welcome to your new account on my service!</h1>
	
		<p>
		  Here is some more information:
	
		  <ul>
			<li>Link 1: <a href="https://example.com">Example.com</a></li>
			<li>Link 2: <a href="https://example2.com">Example2.com</a></li>
			<li>Something else</li>
		  </ul>
		</p>
	
		<table>
		  <thead>
			<tr><th>Header 1</th><th>Header 2</th></tr>
		  </thead>
		  <tfoot>
			<tr><td>Footer 1</td><td>Footer 2</td></tr>
		  </tfoot>
		  <tbody>
			<tr><td>Row 1 Col 1</td><td>Row 1 Col 2</td></tr>
			<tr><td>Row 2 Col 1</td><td>Row 2 Col 2</td></tr>
		  </tbody>
		</table>
	  </body>
	</html>`
	expectedText := `Mega Service ( http://jaytaylor.com/ )

******************************************
Welcome to your new account on my service!
******************************************

Here is some more information:

* Link 1: Example.com ( https://example.com )
* Link 2: Example2.com ( https://example2.com )
* Something else

+-------------+-------------+
|  HEADER 1   |  HEADER 2   |
+-------------+-------------+
| Row 1 Col 1 | Row 1 Col 2 |
| Row 2 Col 1 | Row 2 Col 2 |
+-------------+-------------+
|  FOOTER 1   |  FOOTER 2   |
+-------------+-------------+`
	type args struct {
		n FHIRNarrative
	}
	tests := []struct {
		name string
		args args
		want base.Markdown
	}{
		{
			name: "good case",
			args: args{
				n: FHIRNarrative{
					Div: base.XHTML(inputHTML),
				},
			},
			want: base.Markdown(expectedText),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HumanNarrative(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HumanNarrative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanQuantity(t *testing.T) {
	val := 6.1
	cmp := QuantityComparatorEnumGreaterThan
	unit := "mmol/l"

	type args struct {
		q FHIRQuantity
	}
	tests := []struct {
		name string
		args args
		want base.Markdown
	}{
		{
			name: "good case",
			args: args{
				q: FHIRQuantity{
					Value:      val,
					Comparator: &cmp,
					Unit:       unit,
				},
			},
			want: base.Markdown("> 6.1000 mmol/l"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HumanQuantity(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HumanQuantity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanRange(t *testing.T) {
	unit := "mmol/l"
	highVal := 9.3
	lowVal := 6.1

	type args struct {
		q FHIRRange
	}
	tests := []struct {
		name string
		args args
		want base.Markdown
	}{
		{
			name: "good case",
			args: args{
				q: FHIRRange{
					Low: FHIRQuantity{
						Value:      lowVal,
						Comparator: nil,
						Unit:       unit,
					},
					High: FHIRQuantity{
						Value:      highVal,
						Comparator: nil,
						Unit:       unit,
					},
				},
			},
			want: base.Markdown("6.1000 mmol/l - 9.3000 mmol/l"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HumanRange(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HumanRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanRatio(t *testing.T) {
	denomVal := 9.3
	numerVal := 6.1

	type args struct {
		q FHIRRatio
	}
	tests := []struct {
		name string
		args args
		want base.Markdown
	}{
		{
			name: "good case",
			args: args{
				q: FHIRRatio{
					Numerator: FHIRQuantity{
						Value:      numerVal,
						Comparator: nil,
						Unit:       "",
					},
					Denominator: FHIRQuantity{
						Value:      denomVal,
						Comparator: nil,
						Unit:       "",
					},
				},
			},
			want: base.Markdown("6.1000 / 9.3000"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HumanRatio(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HumanRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}
