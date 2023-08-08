This project was created out of my frustration for having to manually edit an HTML template and manually rendering it into PDF just to issue a simple invoice. I wanted to have a simple program, issue the invoices locally and output them to a directory so that I could back it up easily.

# installing

1. either compile the program from source or grab the binary release from github.
2. install puppeteer core (preferrably, in some other directory)

```sh
npm install puppeteer-core
```

# how does it work?

nothing fancy really. the program takes your invoice definition, then renders it against provided HTML template, then calls puppeteer and saves the output as a PDF file.

# config.yaml

Here's a sample config file:

```yaml
numbering: '{{ year }}/{{ month | stringformat:"%02d" }}/{{ no | stringformat:"%02d" }}'
output: invoices/{{ year }}/{{ invoice }}.pdf
renderer:
  engine: puppeteer
  command:
    - node
    - path-to-pdf.js
    - /usr/bin/brave-browser
    - "{{ source_file }}"
    - "{{ target_pdf }}"
templates: "../invoice-templates"
```

so as you can see, you are defining some generic properties. I'll explain them briefly:

| parameter        | meaning                                                                                                   |
| ---------------- | --------------------------------------------------------------------------------------------------------- |
| numbering        | describes global numbering scheme, for all of the contractors (wihin month)                               |
| output           | describes the template for the target PDF file                                                            |
| renderer         | defines rendering engine. currently the only one supported is puppeteer                                   |
| renderer.command | defines command with it's arguments that will be used to convert a rendered HTML invoice to it's PDF form |
| templates        | specifies the default location for invoice templates                                                      |

note: if you're using MacOS, you may need to replace the path to your browser with a form similar to `/Applications/Brave Browser.app/Contents/MacOS/Brave Browser`

If you will use the default templates location, then the program will look for a file named `<templates-dir>/<contractor-id>/invoice.html`. You can override it within the spec by populating `contractor.template` field.

# invoice-spec.json

```json
{
  "contractor": {
    "id": "contractor-id",
    "pdf-name": "whatever-{{ year }}-{{ month }}-invoice",
    "issue-date": "end-of-month",
    "rate": [
      { "start-date": "2020-01-01", "rate": 1 },
      { "start-date": "2023-01-01", "rate": 2 }
    ]
  },
  "items": [
    {
      "name": "First item of the invoice",
      "quantity": "month-workdays",
      "days_free": 0
    }
  ]
}
```

With the above example, the program will:

1. calculate the number of workdays
2. subtract number of free days
3. select appropriate rate and generate the structure
4. set the issue date as the last day of the previous month. In other words, if you run the command on the 1'st day of July, the invoice would be generated for the month of June, on the last day of it.
5. save it as `invoices/(year)/whatever-(year)-(month)-invoice.pdf

# invoice-with-vat

```json
{
  "contractor": {
    "id": "contractor-id",
    "pdf-name": "whatever-{{ year }}-{{ month }}-invoice",
    "issue-date": "end-of-month",
    "rate": [
      {
        "start-date": "2020-01-01",
        "rate": 1,
        "gross": true,
        "vat": { "rate": 0.23 }
      },
      { "start-date": "2023-01-01", "rate": 2, "vat": { "rate": 0.23 } }
    ]
  },
  "items": [
    {
      "name": "First item of the invoice",
      "quantity": "month-workdays",
      "days_free": 0
    },
    {
      "name": "additional item",
      "quantity": 1,
      "unit_price": {
        "price": 100,
        "gross": true
      },
      "vat": {
        "rate": 0.08
      }
    }
  ]
}
```

in this case, the program will calculate the net and gross amounts based on VAT rate. Also, if you won't provide a label to the VAT rate, the default label would be rendered as rate \* 100 + % (therefore 0.23 would become 23 %). a mapping of totals per vat rate will also be generated

# rendering options

you can render the invoice to PDF (default), to HTML (I suppose for checking in the browser) or dump the structure as JSON. for details, please use `./kwity issue --help`.

If you re-render the same invoice within the same month, then the number will not increase. instead, the invoice will simply be re-rendered.

# rendering script for puppeteer

it's a very slightly modified example from puppeteer website (please refer to examples/pdf.js)
