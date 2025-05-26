<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/context7-http?layout=left&icon=fluent-emoji-flat%3Amagic-wand&icon.width=60&bg=geometric)

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">
  <a href="https://github.com/lrstanley/context7-http/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/context7-http?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/context7-http/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/context7-http?style=flat-square">
  </a>




  <a href="https://github.com/lrstanley/context7-http/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/actions/workflow/status/lrstanley/context7-http/test.yml?branch=master&label=test&style=flat-square">
  </a>


  <a href="https://codecov.io/gh/lrstanley/context7-http">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/context7-http/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/context7-http">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/context7-http?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/context7-http">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/context7-http?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/context7-http/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/context7-http/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/context7-http/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/context7-http/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/context7-http/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/context7-http?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/context7-http/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Features](#sparkles-features)
  - [Usage](#gear-usage)
    - [VSCode, Cursor, etc](#vscode-cursor-etc)
    - [Install in Windsurf](#install-in-windsurf)
    - [Install in Zed](#install-in-zed)
    - [Install in Claude Code](#install-in-claude-code)
    - [Install in Claude Desktop](#install-in-claude-desktop)
    - [Install in BoltAI](#install-in-boltai)
  - [References](#books-references)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
<!-- template:end:toc -->

## :sparkles: Features

**context7-http** is a MCP server that supports HTTP streaming for the [Context7](https://context7.com) project.
This allows you to utilize the MCP server from anywhere, without installing anything locally.

- Has _current_ feature parity with the existing Context7 MCP Server.
- SSE and HTTP `streamable` support.
- Provides `resolve-library-uri` and `search-library-docs` tools for finding libraries, and searching their documentation.
- Provides multiple resources, including:
  - `context7://libraries` - returns high-level information about all libraries.
  - `context7://libraries/<project>`
  - `context7://libraries/top/<n>` - returns the top `n` libraries, sorted by trust score (if available), otherwise by stars.

---

## :gear: Usage

For all examples below, replace `context7.liam.sh` with your own MCP server URL, if you're running your own instance.

Configured endpoints:

- `https://context7.liam.sh/mcp` - HTTP `streamable` endpoint.
- `https://context7.liam.sh/sse` (& `/message`) - SSE endpoint (**NOTE**: SSE is considered deprecated in the MCP spec).

Other than swapping out the `mcpServer` block (or similar, depending on your client), usage should match that of the
[official Context7 documentation](https://github.com/upstash/context7#-with-context7)

### VSCode, Cursor, etc

[Cursor MCP docs](https://docs.cursor.com/context/model-context-protocol#configuring-mcp-servers), and
[VS Code MCP docs](https://code.visualstudio.com/docs/copilot/chat/mcp-servers) for more info.

```json
{
  "mcpServers": {
    "context7": {
      "url": "https://context7.liam.sh/mcp"
    }
  }
}
```

### Install in Windsurf

Add this to your Windsurf MCP config file. See [Windsurf MCP docs](https://docs.windsurf.com/windsurf/cascade/mcp#mcp-config-json) for more info.

```json
{
  "mcpServers": {
    "context7": {
      "url": "https://context7.liam.sh/mcp"
    }
  }
}
```

### Install in Zed

Add this to your Zed `settings.json`. See [Zed MCP docs](https://zed.dev/docs/ai/mcp#bring-your-own-mcp-server) for more info.

```json
{
  "context_servers": {
    "context7": {
      "url": "https://context7.liam.sh/mcp"
    }
  }
}
```

### Install in Claude Code

Run this command. See [Claude Code MCP docs](https://docs.anthropic.com/en/docs/claude-code/tutorials#configure-mcp-servers) for more info.

```sh
claude mcp add --transport sse context7 https://context7.liam.sh/mcp
```

### Install in Claude Desktop

Add this to your Claude Desktop `claude_desktop_config.json` file. See [Claude Desktop MCP docs](https://modelcontextprotocol.io/quickstart/user) for more info.

```json
{
  "mcpServers": {
    "context7": {
      "url": "https://context7.liam.sh/mcp"
    }
  }
}
```

### Install in BoltAI

[BoltAI MCP docs](https://docs.boltai.com/docs/plugins/mcp-servers#how-to-use-an-mcp-server-in-boltai).

## :books: References

- [Context7](https://context7.com) - [repo](https://github.com/upstash/context7)
- [Model Context Protocol Introduction](https://modelcontextprotocol.io/introduction)

---

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
* :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
* :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/context7-http/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
    community.
* :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
* :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/context7-http/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2025 Liam Stanley <liam@liam.sh>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->
