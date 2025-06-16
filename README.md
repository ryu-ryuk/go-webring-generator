# Go Webring Generator

A simple command-line tool written in Go to generate static files for a webring, including navigation links and an index page.

## Features

* Generates `index.html` with a list of webring members.
* Creates `/rand/index.html` for a random redirect.
* Generates `/slug/previous/index.html` and `/slug/next/index.html` for each member, redirecting to the previous and next sites in the ring.
* Uses TOML for data configuration.
* Supports custom HTML templates.

## Getting Started

### Prerequisites

* Go 1.16 or higher

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/ryu-ryuk/go-webring-generator.git
    cd go-webring-generator
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

### Usage

1.  **Configure Members:**
    Edit `data/members.toml` to define your webring members.

    Example `data/members.toml`:
    ```toml
    [[members]]
    name = "My Awesome Site"
    slug = "my-awesome-site"
    about = "A site about cool stuff."
    url = "[https://example.com/my-awesome-site](https://example.com/my-awesome-site)"
    github = "your-github-username"
    owner = "Your Name"
    role = "Developer"

    [[members]]
    name = "Another Great Blog"
    slug = "another-blog"
    about = "Thoughts on technology and life."
    url = "[https://anotherblog.net](https://anotherblog.net)"
    github = "another-github-user"
    owner = "Another Person"
    role = "Writer"
    ```

2.  **Customize Templates (Optional):**
    Modify the HTML templates in the `templates/` directory to match your desired layout and styling.

    * `templates/index.html`
    * `templates/random.html`
    * `templates/redirect.html`

3.  **Add Static Assets (Optional):**
    Place any static files like CSS (`styles.css`) in the `static/` directory.

4.  **Generate Static Files:**
    Run the generator:
    ```bash
    go run main.go
    ```
    This will create a `public/` directory containing all the generated HTML files and copied static assets.

## Project Structure
```
go-webring-generator/
├── main.go # main Go application logic
├── go.mod # Go module definition
├── go.sum # Go module checksums
├── data/
│   └── members.toml # Webring member data in TOML format
├── templates/
│   ├── index.html # Template for the main index page
│   ├── random.html # Template for the random redirect page
│   └── redirect.html # Template for the previous/next redirects
└── static/
    └── styles.css # Static CSS file 
```

## Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
