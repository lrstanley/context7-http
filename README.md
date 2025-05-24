<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/context7-http?layout=left&icon=fluent-emoji-flat%3Amagic-wand&icon.width=60&bg=geometric)

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Features](#sparkles-features)
  - [Usage](#gear-usage)
  - [References](#books-references)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
<!-- template:end:toc -->

## :sparkles: Features

**context7-http** is a MCP server that supports HTTP SSE streaming for the [Context7](https://context7.com) project.

- HTTP SSE streaming for Context7.
- (_future_) Support for HTTP `streamable` functionality once it's supported in the upstream MCP library
  we use (or the official MCP library is released by the Go team).
- Provides `resolve-library-uri` and `search-library-docs` tools for finding libraries, and searching their documentation.
- Provides multiple resources, including:
  - `context7://libraries` - returns high-level information about all libraries.
  - `context7://libraries/<project>`
  - `context7://libraries/top/<n>` - returns the top `n` libraries, sorted by trust score (if available), otherwise by stars.

---

## :gear: Usage

TODO

## :books: References

- [Context7](https://context7.com) - [repo](https://github.com/upstash/context7)
- [Model Context Protocol Introduction](https://modelcontextprotocol.io/introduction)

---

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
<!-- template:end:license -->
