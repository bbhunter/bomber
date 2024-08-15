// Package ai contains functionality to render output using GenAI
package ai

import (
	"bytes"
	"html/template"
	"log"

	"github.com/devops-kung-fu/common/util"
	"github.com/spf13/afero"

	"github.com/devops-kung-fu/bomber/lib"
	"github.com/devops-kung-fu/bomber/models"
)

// Renderer contains methods to render AI HTML output format
type Renderer struct{}

// Render outputs ai generated report
func (Renderer) Render(results models.Results) error {
	var afs *afero.Afero

	lib.MarkdownToHTML(results) 

	if results.Meta.Provider == "test" { 
		afs = &afero.Afero{Fs: afero.NewMemMapFs()}
	} else {
		afs = &afero.Afero{Fs: afero.NewOsFs()}
	}

	filename := lib.GenerateFilename()
	util.PrintInfo("Writing AI Enriched HTML report:", filename)

	resultString, err := generateTemplateResult(templateString(), results)
	if err != nil {
		log.Println(err)
		return err
	}
	err = afs.WriteFile(filename, []byte(resultString), 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func generateTemplateResult(templateString string, data interface{}) (string, error) {
	// Create a new template with a name
	tmpl, err := template.New("output").Parse(templateString)
	if err != nil {
		return "", err
	}

	// Create a buffer to store the generated result
	var resultBuffer bytes.Buffer

	// Execute the template and write the result to the buffer
	err = executeTemplate(&resultBuffer, tmpl, data)
	if err != nil {
		return "", err
	}

	// Convert the buffer to a string and return it
	return resultBuffer.String(), nil
}

func executeTemplate(buffer *bytes.Buffer, tmpl *template.Template, data interface{}) error {
	// Execute the template and write the result to the buffer
	return tmpl.Execute(buffer, data)
}

func templateString() string {
	return `

	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>bomber Results</title>
	
		<style>
		body {
			font-family: Helvetica;
			margin: 20px;
		}
		#vuln {
			border: 1px solid;
			border-color: gray;
			border-radius: 5px;
			padding: 10px;
			margin-bottom: 10px;
		}
		#bomber-logo {
			margin-top: 10px;
			margin-bottom: 30px;
			width:331px;
			height:80px;
			background-image:url(data:image/png;base64,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);
		}
		#summary {
			border: 0px;
		}
		#CRITICAL {
			color: red;
			font-weight: bold;
		}
		#HIGH {
			color: purple;
			font-weight: bold;
		}
		#MODERATE {
			color: goldenrod;
			font-weight: bold;
		}
		#LOW {
			color: green;
			font-weight: bold;
		}
		#UNSPECIFIED {
			color: blue;
			font-weight: bold;
		}
		h3 {
			font-weight: bold
		}
		#footer {
			width: 75px;
			height: 75px;
			margin-bottom: 10px;
			background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEsAAABLCAYAAAA4TnrqAAAAAXNSR0IArs4c6QAAAIRlWElmTU0AKgAAAAgABQESAAMAAAABAAEAAAEaAAUAAAABAAAASgEbAAUAAAABAAAAUgEoAAMAAAABAAIAAIdpAAQAAAABAAAAWgAAAAAAAABIAAAAAQAAAEgAAAABAAOgAQADAAAAAQABAACgAgAEAAAAAQAAAEugAwAEAAAAAQAAAEsAAAAAdzEKuAAAAAlwSFlzAAALEwAACxMBAJqcGAAAAVlpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IlhNUCBDb3JlIDYuMC4wIj4KICAgPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICAgICAgPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIKICAgICAgICAgICAgeG1sbnM6dGlmZj0iaHR0cDovL25zLmFkb2JlLmNvbS90aWZmLzEuMC8iPgogICAgICAgICA8dGlmZjpPcmllbnRhdGlvbj4xPC90aWZmOk9yaWVudGF0aW9uPgogICAgICA8L3JkZjpEZXNjcmlwdGlvbj4KICAgPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KGV7hBwAAGR1JREFUeAHVnAWUZcW1hs8M7u4QXIIEXhjkwePh7rZIkAQnggV3XQt3dx8cgjMwDO42WAgEGSS4+3T37a73f7vrP6m+dM/0zDQwb69Vt86pU7Jr17baVd39qp8extUQ/ZXalFIx3Dh6nlZpRqUZcppZ+UxKUyl9p/RhTp/k/DPlHyt9qVTCeHqh70ZZ2NfP/fq6w9wf/UKMZuTnHGeccZafcsop151gggkWHm+88aZIKU3fT9C/f/9K36q2trbqyy+/rFSn0veqo6MjkurR9feNRuOLlpaWoV999dUdqnuvyl7lQwEsTrtSuTDF59F/7GtiwUEQCS4yLK6HVaabbro1p5pqqgGiyyRvvPFGpUn7O3mHEpNzThkAfk70XeM777zzViJwiwg79KOPPrpL3+5ReljJALdBNPrsE6gH74PeQM5EYmLbKu0wyyyzLAmX/OMf//AQMYGNNtqo+u1vf9tv4okn7q/Ub4YZZug300wzVRNNNFGkH374oWptba0+/vjj6tNPP63ESam9vb3jlVdeSeeeey6E9cJUv/71r6tvv/22evfdd59W+QlK13gw5SVeRfEv8wgnkQxbSsSem3322ZMIxaRIELGxxx57tP/973/veP3115MmJ8kadZAIpjfffDPdfvvtHfvvvz+Eh0Xpv32OOeZIv/rVr5KI/6TeN1My1IR1wejkY8pZ6AfL0//q+ei55pprGUTsnXfeYQL9/vKXv4y72mqrwUWViIfo1HiKU0Ic0Uuff/453FMNHz68EgkriWsloleTTjpp6K8JJ5wwysYff/y6PQ8ffvhhcO3NN99cnX766YGLxhkXDtWiPKEqByjdS13BGHHZ6BKLdqHApYQXlaI9WKK28cwzz1y9/PLLIWYHHHDAeL///e9DRFDchm+++aYaNmxY9dxzz1VPPPFETPbrr7+uHn744SCU65X5MsssU00xxRTV1FNPXa244orVb37zm2q++eaLMtdjgV588cXquuuuq44++mgWqv/8888/DuO9//77g0Xk3SXWL6scgkFUOP4nB1jD7LGFnpO4iYEhUpsI1PHUU08lcUstY4jc0KFDk1Y+rb766hbNH+XSW2nOOedMc889d+TSYUmc+KN6jLnccsulE088MT3++ONJxK7HYlzG+tOf/oRih2iNeeaZhz5aldZXArzYnW8/0S+rYth3kkkmSQsssACItCy77LLpjjvuSN9//32N+BdffJGuvvrqtOaaa3aZMPqMduLENM000yRxTRKHdqmjPuNdXJkmm2yyJM5NsqjRLi9OXX+NNdZIN910U2I8A7rtscceS1tttRX1Whhv8skn5/lAJUM5H5f1SY5+MlwikUizzTZbrN4ZZ5yRxO7GM8mSJemQNGDAgHpCcA0Iq02P3CI9lZqTBqz7KJ9nnXXWJDELhe5y6cV04403djEe0oG2nm0zzjhjB9yq+jcvvvjiJlQ5L30ac4gOp5122snU1W1YHeWweLt0RC1y0l1JeijtuOOO9QQRgemnnx4LVZepXRCFfFQTxCzbyAgEx2VRi2+bbLJJuvvuuxPcZbj11lv5hqpozXWHymjMpXfAhOt8G4PfIJSs0nTq45UsAi16DsIYGVmzdOihh9YTWXDBBUO8qNc8Qcp6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);}
		</style>
	</head>
	<body>
		<div id="bomber-logo"></div>
		<p>The following results were detected by <code>{{.Meta.Generator}} {{.Meta.Version}}</code> on {{.Meta.Date}} using the {{.Meta.Provider}} provider.</p>
		{{ if ne (len .Packages) 0 }} 
		<p>
			Vulnerabilities displayed may differ from provider to provider. This list may not contain all possible vulnerabilities. Please try the other providers that <code>bomber</code> supports (osv, ossindex, snyk). There is no guarantee that the next time you scan for vulnerabilities that there won't be more, or less of them. Threats are continuous.
		</p>
		<p>
			<b>This report contains vulnerability explanations that have been generated by AI. All explanations should be validated and may differ between scans.</b>
		</p>
		<p>
			EPSS Percentage indicates the % chance that the vulnerability will be exploited. This
			value will assist in prioritizing remediation. For more information on EPSS, refer to
			<a href="https://www.first.org/epss/">https://www.first.org/epss/"</a>.</p>
		</p>
		{{ else }}
		<p>
			No vulnerabilities found!
		</p>
		{{ end }}
		{{ if ne (len .Files) 0 }} 
			<h1>Scanned Files</h1>
			{{ range .Files }}
				<p><b>{{ .Name }}</b> (sha256:{{ .SHA256 }})</p>
			{{ end }}
		{{end}}
		{{ if ne (len .Licenses) 0 }} 
			<h1>Licenses</h1>
			<p>The following licenses were found by <code>bomber</code>:</p>
			<ul>
			{{ range $license := .Licenses }}
				<li>{{ $license }}</li>
			{{ end }}
			</ul>
		{{ else }}
			<p>No license information detected.</b>
		{{ end }}
		{{ if ne (len .Packages) 0 }} 
			<h1>Vulnerability Summary</h1>
			{{ if ne (len .Meta.SeverityFilter) 0 }}
				<p>Only showing vulnerabilities with a severity of <i><b>{{ .Meta.SeverityFilter }}</b></i> or higher.</p>
			{{ end }}
			<table id="summary">
				{{if gt .Summary.Critical 0}}
				<tr><td>Critical:</td><td>{{ .Summary.Critical }}</td></tr>
				{{ end }}
				{{if gt .Summary.High 0}}
				<tr><td>High:</td><td>{{ .Summary.High }}</td></tr>
				{{ end }}
				{{if gt .Summary.Moderate 0}}
				<tr><td>Moderate:</td><td>{{ .Summary.Moderate }}</td></tr>
				{{ end }}
				{{if gt .Summary.Low 0}}
				<tr><td>Low:</td><td>{{ .Summary.Low }}</td></tr>
				{{ end }}
				{{if gt .Summary.Unspecified 0}}
				<tr><td>Unspecified:</td><td>{{ .Summary.Unspecified }}</td></tr>
				{{ end }}
			</table>
			<h1>Vulnerability Details</h1>
			{{ range .Packages }}
				<h2>{{ .Purl }}</h2>
				<p>{{ .Description }}</p>
				<h3>Vulnerabilities</h3>
				{{ range .Vulnerabilities }}
					<div id="vuln">
						{{ if .Title }}
						<h3>{{ .Title }}</h3>
						{{ end }}
						<p>Severity: <span id="{{ .Severity }}">{{ .Severity }}</span></p>
						{{ if ne (len .Epss.Percentile) 0 }} 
							<p>EPSS: <span>{{ .Epss.Percentile }}</span></p>
						{{ end }}
						{{ if .Explanation }}
							<p><b>Vulnerability Explanation</b></p>
							<p>{{ .Explanation }}</p>
						{{ end }}
						<p><a href="{{ .Reference }}">Reference Documentation</a></p>
						<p><b>Details</b></p>
						{{ .Description }}
					</div>
				{{ end }}
				<br/>
			{{ end }}
		{{ end }}
		<div id="footer"></div>
		Powered by the <a href="https://github.com/devops-kung-fu"/>DevOps Kung Fu Mafia</a>
	</body>
	</html>
	
	`
}