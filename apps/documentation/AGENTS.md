# AGENTS.md: Documentation Site

> **AUTHORITY:** This document governs AI agent operations within `apps/documentation/`.
> **INHERITANCE:** Inherits from `/AGENTS.md` and extends with documentation-specific rules.
> **PRECEDENCE:** Rules here OVERRIDE root rules for this directory.

## Purpose

Defines strict operational rules for agents working on documentation. Ensures consistency with:

- Google Developer Documentation Style Guide
- Technical writing best practices
- Documentation structure and organization
- Content quality standards

## Scope

This file governs:

- `apps/documentation/pages/` - Documentation pages
- `apps/documentation/docs/` - ADRs and runbooks
- `apps/documentation/src/` - Custom components and pages
- `docusaurus.config.ts` - Site configuration
- `sidebars.ts` - Navigation structure
- `package.json` - Dependencies and scripts
- Style and formatting standards

## Mandatory Workflow

**MUST follow planning-with-files workflow from `/.github/skills/planning-with-files/`**

This is NON-NEGOTIABLE for:

- Creating new documentation sections
- Major content restructuring
- Style guide changes
- Multi-page documentation updates

See root `/AGENTS.md` for workflow details.

## Documentation Goals

Documentation MUST be:

- **Clear** - Easy to understand on first read
- **Precise** - Technically accurate and complete
- **Boring** - Consistent and predictable
- **Scannable** - Easy to skim and navigate
- **Maintainable** - Easy to update and extend

Documentation is NOT:

- Marketing copy or promotional content
- Blog posts or opinion pieces
- Creative writing exercises
- Entertainment

**Documentation is a contract with the reader.**

## Style Guide Authority

**Primary Reference:** [Google Developer Documentation Style Guide](https://developers.google.com/style)

All documentation MUST conform to this guide unless explicitly overridden below.

**Key Principles:**

1. **Voice and Tone**

   - Second person ("you")
   - Present tense
   - Active voice
   - Neutral and professional

2. **Clarity**

   - Short sentences
   - One idea per paragraph
   - Simple words over complex ones
   - No jargon without definition

3. **Consistency**

   - Same term for same concept
   - Parallel grammar in lists
   - Predictable structure

4. **Precision**
   - No fluff or filler
   - No hype or marketing language
   - No anthropomorphism
   - No jokes or metaphors

**If it sounds clever, rewrite it.**

## Document Structure

### Required Structure

Every documentation page MUST follow this order (unless exempted):

1. **Title** (H1)
2. **Summary paragraph** (1-2 sentences describing the page)
3. **Prerequisites** (if applicable)
4. **Main content** (organized with clear headings)
5. **Examples** (if applicable)
6. **Next steps** (optional)

### Title Rules

**Format:**

- Sentence case (not title case)
- Descriptive, not clever
- No punctuation at end
- No emojis or symbols

**Examples:**

✅ Good:

- Configure authentication
- Install dependencies
- Database migration guide

❌ Bad:

- 🔐 Auth Setup Made Easy!
- How To Install Dependencies
- Database Migrations.

### Heading Hierarchy

**Rules:**

- One H1 per page (the title)
- H2 for main sections
- H3 for subsections
- H4 for sub-subsections
- Never skip levels (H2 → H4 is invalid)
- Use sentence case

**Example:**

```markdown
# Configure authentication

## Prerequisites

### Required tools

## Setup process

### Install packages

### Configure environment

## Verify installation
```

### Paragraphs

**Rules:**

- One idea per paragraph
- 2-4 sentences maximum
- No wall-of-text paragraphs
- Use line breaks between paragraphs

### Lists

**Bulleted Lists:**

- Use for unordered collections
- Use for items without sequence
- Parallel grammar in all items
- Complete sentences end with periods
- Fragments don't use periods

**Numbered Lists:**

- Use for sequences and steps
- Use for prioritized items
- Start each item with verb (for steps)
- Parallel grammar in all items

**Examples:**

✅ Good (bulleted):

```markdown
Features include:

- Fast build times
- Hot module replacement
- TypeScript support
```

✅ Good (numbered):

```markdown
To deploy:

1. Build the application
2. Run tests
3. Push to production
4. Verify deployment
```

❌ Bad (inconsistent grammar):

```markdown
To deploy:

1. Building the application
2. Run tests
3. You should push to production
```

## Writing Rules

### Voice and Tone

**Use second person:**

✅ Correct:

> Run the command to generate the build.

❌ Incorrect:

> We can run the command to generate the build.
> The user should run the command to generate the build.

**Use active voice:**

✅ Correct:

> The system validates your credentials.

❌ Incorrect:

> Your credentials are validated by the system.

**Use imperative verbs for instructions:**

✅ Correct:

> Install the package using pnpm.

❌ Incorrect:

> You should install the package using pnpm.
> The package should be installed using pnpm.

**Use present tense:**

✅ Correct:

> The function returns a promise.

❌ Incorrect:

> The function will return a promise.
> The function returned a promise.

### Capitalization

**Sentence case for headings:**

✅ Correct:

> ## Configure the database

❌ Incorrect:

> ## Configure The Database
>
> ## Configure the Database

**Proper nouns:**

- Capitalize: Next.js, TypeScript, Supabase, GitHub
- Don't capitalize: website, database, server, client

**Features and concepts:**

- Don't capitalize unless branded: authentication (not Authentication)
- Exception: When referring to UI elements ("Click the Authentication button")

### Terminology

**Consistency rules:**

- Use ONE term per concept throughout all documentation
- Do not introduce synonyms
- Define terms on first use
- Create glossary entry for domain-specific terms

**Examples:**

✅ Consistent:

> The **authentication system** validates user credentials. Configure authentication in the settings.

❌ Inconsistent:

> The **auth system** validates user credentials. Configure **user verification** in the settings.

**Common terms:**

| Use This       | Not This                                    |
| -------------- | ------------------------------------------- |
| authentication | auth, login system                          |
| command        | instruction, directive                      |
| repository     | repo, codebase                              |
| function       | method (unless specifically a class method) |

### Warnings and Notes

**Use these labels consistently:**

**Note:** Adds context or clarification.

```markdown
**Note:** This feature requires Node.js 18 or later.
```

**Caution:** Indicates potential problems.

```markdown
**Caution:** This operation cannot be undone.
```

**Warning:** Indicates risk of data loss or system failure.

```markdown
**Warning:** This will delete all existing data.
```

**Rules:**

- Do not overuse (max 1-2 per page)
- Place before the action, not after
- Be specific about the risk or context

### Forbidden Patterns

**MUST NOT use:**

❌ Emojis (unless in code examples that use emojis)

❌ Anthropomorphism:

> ❌ The system thinks the input is invalid
> ✅ The system determines that the input is invalid

❌ Marketing language:

> ❌ Our amazing authentication system
> ✅ The authentication system

❌ Future tense:

> ❌ The function will process your request
> ✅ The function processes your request

❌ Passive voice:

> ❌ The request is processed by the server
> ✅ The server processes the request

❌ Jokes or humor:

> ❌ Let's dive into the database (unless you have scuba gear!)
> ✅ Configure the database connection

❌ Metaphors:

> ❌ The data flows like a river through the pipeline
> ✅ The system transfers data through the pipeline

❌ Exclamation marks (except in warnings)

## Code Documentation

### Code Blocks

**Rules:**

- MUST specify language
- Keep examples minimal
- Make examples runnable
- Do not use placeholder code without marking it

**Format:**

````markdown
```typescript
// Good: Specific and runnable
const result = await fetchUser("123");
```

```typescript
// Bad: Vague placeholder
const result = await fetchData(/* your data here */);
```
````

### Inline Code

**Use inline code for:**

- Commands: `pnpm install`
- Filenames: `package.json`
- Function names: `fetchUser()`
- Variable names: `userId`
- Package names: `next`
- Flags: `--filter`
- Keyboard keys: `Ctrl+C`

**Do not use for:**

- Emphasis (use **bold**)
- File paths longer than one segment (use links)

### Command Examples

**Format:**

````markdown
Install dependencies:

```bash
pnpm install
```
````

Start the server:

```bash
mise server:start
```

````

**Rules:**
- Show full command, not partial
- Include required flags
- Explain output if relevant
- Use `bash` or `sh` as language

### API Documentation

**Function documentation:**

```markdown
## fetchUser

Retrieves user data by ID.

**Parameters:**
- `id` (string): The user ID
- `options` (FetchOptions, optional): Request options

**Returns:**
- Promise<User>: User object

**Throws:**
- `UserNotFoundError`: When user doesn't exist
- `NetworkError`: When request fails

**Example:**

```typescript
const user = await fetchUser('123', { cache: 'no-store' })
console.log(user.name)
````

````

## MDX and React Components

Docusaurus uses MDX, allowing React components within Markdown. Use this for enhanced documentation with interactive elements.

### MDX Usage Rules

**When to use MDX:**

- Interactive code examples
- Custom callouts or warnings
- Reusable UI components
- Complex diagrams or visualizations

**Import rules:**

- Import components at the top of the file
- Use relative paths for local components
- Group imports logically

**Example:**

```mdx
import BrowserWindow from '@site/src/components/BrowserWindow';
import CodeBlock from '@theme/CodeBlock';

<BrowserWindow>
  <CodeBlock language="js">
    {`console.log('Hello, world!');`}
  </CodeBlock>
</BrowserWindow>
````

**JSX Guidelines:**

- Keep JSX simple and readable
- Use descriptive prop names
- Avoid complex logic in MDX files
- Test components in isolation

### Custom Components

**Location:** `src/components/`

**Rules:**

- Use TypeScript for type safety
- Document component props with JSDoc
- Follow React best practices
- Keep components focused and reusable

**Example Component:**

```typescript
// src/components/Alert.tsx
interface AlertProps {
  type: "info" | "warning" | "error";
  children: React.ReactNode;
}

/**
 * Displays an alert message
 */
export default function Alert({ type, children }: AlertProps) {
  return <div className={`alert alert--${type}`}>{children}</div>;
}
```

**Usage in MDX:**

```mdx
import Alert from "@site/src/components/Alert";

<Alert type="warning">This is a warning message.</Alert>
```

## File Organization

### Directory Structure

```
apps/documentation/
├── docs/                         # Documentation pages
│   ├── intro.md                  # Introduction
│   ├── company/                  # Company documentation
│   ├── adr/                      # Architecture Decision Records
│   └── runbooks/                 # Operational procedures
├── pages/                        # Standalone pages
│   ├── index.md                  # Home page
│   ├── guide/                    # User guides
│   ├── api/                      # API reference
│   └── tutorials/                # Step-by-step tutorials
└── src/                          # Custom components and pages
    ├── components/               # Reusable React components
    │   └── HomepageFeatures/     # Homepage components
    ├── css/                      # Custom styles
    └── pages/                    # Custom React pages
```

**Root Configuration Files:**

- `docusaurus.config.ts` - Site configuration
- `sidebars.ts` - Navigation structure
- `package.json` - Dependencies and scripts
- `tsconfig.json` - TypeScript configuration

### File Naming

**Rules:**

- Lowercase only
- Use hyphens for spaces: `getting-started.md`
- Descriptive names: `database-migrations.md` (not `db-mig.md`)
- No dates in filenames (use frontmatter)

### Frontmatter

**Required fields:**

```markdown
---
title: Page title
description: Brief description for SEO
---
```

**Optional fields:**

```markdown
---
title: Page title
description: Brief description
author: Author name
date: 2026-01-07
tags: [tag1, tag2]
---
```

## Cross-References

### Internal Links

**Format:**

```markdown
See [Getting Started](./getting-started.md) for more information.

Configure [authentication settings](../guide/authentication.md).
```

**Rules:**

- Use relative paths
- Include `.md` extension
- Use descriptive link text (not "click here")

### External Links

**Format:**

```markdown
See the [Next.js documentation](https://nextjs.org/docs) for details.
```

**Rules:**

- Use full URLs
- Open in same window (no `target="_blank"` in Markdown)
- Verify links are current

## Architecture Decision Records (ADRs)

**Location:** `apps/documentation/docs/adr/`

**Naming:** `NNNN-title.md` (e.g., `0001-use-typescript.md`)

**Template:**

```markdown
# ADR-NNNN: Title

**Status:** Proposed | Accepted | Deprecated | Superseded

**Date:** YYYY-MM-DD

**Authors:** Names

## Context

What is the issue we're facing?

## Decision

What decision did we make?

## Consequences

What are the positive and negative outcomes?

## Alternatives Considered

What other options did we evaluate?
```

## Runbooks

**Location:** `apps/documentation/docs/runbooks/`

**Purpose:** Operational procedures and troubleshooting guides

**Template:**

````markdown
# Runbook: Title

## Purpose

What does this runbook cover?

## Prerequisites

What you need before starting.

## Steps

1. First step
2. Second step
3. Third step

## Troubleshooting

### Issue: Description

**Symptoms:**

- Symptom 1
- Symptom 2

**Solution:**

1. Solution step 1
2. Solution step 2

## References

- [Link to related docs]

## Configuration Management

Docusaurus configuration requires careful management to maintain site stability and performance.

**Blog Feature Removed:** The blog feature has been disabled to focus on documentation-only content. Blog configuration should not be re-enabled without governance approval.

### docusaurus.config.ts

**Rules:**

- Use TypeScript for type safety
- Document plugin configurations
- Test configuration changes locally
- Follow semantic versioning for config updates

**Common configurations:**

```typescript
// docusaurus.config.ts
export default {
  title: "Documentation",
  tagline: "Technical documentation",
  url: "https://docs.example.com",
  baseUrl: "/",
  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: "./sidebars.ts",
          editUrl: "https://github.com/...",
        },
        blog: {
          showReadingTime: true,
          editUrl: "https://github.com/...",
        },
        theme: {
          customCss: "./src/css/custom.css",
        },
      },
    ],
  ],
};
```
````

### sidebars.ts

**Rules:**

- Use descriptive labels
- Group related pages logically
- Keep nesting shallow (max 3 levels)
- Update when adding/removing pages

**Example:**

```typescript
// sidebars.ts
export default {
  docs: [
    "intro",
    {
      type: "category",
      label: "Company",
      items: ["company/mission", "company/values"],
    },
  ],
};
```

### Plugin Management

**Rules:**

- Only add plugins with clear justification
- Document plugin purpose and configuration
- Test plugins thoroughly before deployment
- Keep plugins updated

**Common plugins:**

- `@docusaurus/plugin-content-docs` - Documentation pages
- `@docusaurus/plugin-content-blog` - Blog functionality
- `@docusaurus/plugin-sitemap` - SEO sitemap
- `@docusaurus/plugin-google-analytics` - Analytics

## Task Execution

**Available mise tasks for documentation:**

```bash
mise docs:setup     # Set up documentation site
mise docs:dev       # Start dev server
mise docs:build     # Build static site
mise docs:lint      # Check for broken links and style
```

**Pre-execution:**

1. Verify task exists: `mise tasks | grep docs`
2. Check for errors in output

**Post-execution:**

1. Verify build success
2. Check for broken links
3. Review changes in browser

## Allowed Actions

Agents working in `apps/documentation/` MAY:

- Create or modify documentation pages
- Create or update ADRs
- Create or update runbooks
- Create or modify custom React components
- Update docusaurus.config.ts with justification
- Modify sidebars.ts for navigation changes
- Add or update dependencies in package.json
- Fix broken links
- Improve clarity and readability
- Add code examples
- Update outdated information
- Reorganize content for better navigation
- Add cross-references

## Forbidden Actions

Agents working in `apps/documentation/` MUST NOT:

- Modify files outside `apps/documentation/` without cross-cutting plan
- Skip the planning workflow for complex tasks
- Modify docusaurus.config.ts without testing and justification
- Break existing navigation in sidebars.ts
- Add untested plugins or dependencies
- Use marketing language or hype
- Add emojis or creative formatting
- Use jokes or metaphors
- Write in first person
- Use passive voice
- Skip required structure elements
- Create unclear headings
- Use title case for headings
- Introduce synonyms for existing terms
- Add content without verifying accuracy
- Copy-paste from marketing materials

## Quality Checklist

Before completing documentation work, verify:

- [ ] Follows Google Style Guide
- [ ] Uses second person and active voice
- [ ] Uses sentence case for headings
- [ ] Has required structure (title, summary, etc.)
- [ ] Code examples are runnable
- [ ] MDX components import correctly
- [ ] Custom components are properly typed
- [ ] Links work correctly
- [ ] Navigation updates reflect in sidebars.ts
- [ ] Configuration changes don't break build
- [ ] No typos or grammar errors
- [ ] Consistent terminology throughout
- [ ] No marketing language or hype
- [ ] No emojis or jokes
- [ ] Clear and scannable
- [ ] Technically accurate

## Review Process

**Self-review questions:**

1. Can a first-time reader understand this?
2. Is every sentence necessary?
3. Is the voice consistent (second person, active)?
4. Are code examples runnable?
5. Is terminology consistent with other docs?
6. Are headings descriptive and in sentence case?
7. Does it follow the required structure?

## Output Expectations

**When completing work:**

1. State: "Working under authority of `apps/documentation/AGENTS.md`"
2. List modified files with absolute paths
3. Confirm no broken links: `mise docs:lint`
4. Confirm build succeeds: `mise docs:build`
5. Provide summary of changes
6. Update `task_plan.md` with completion status

**When reporting errors:**

1. Log to `task_plan.md` under "Errors Encountered"
2. Include full error message
3. Include file path and line number
4. Propose resolution
5. Execute resolution or request guidance

## Maintenance

Update this file when:

- Documentation platform changes
- Style guide requirements updated
- New documentation types added
- Quality standards change

**Update process:**

1. Create plan for governance change
2. Update this file
3. Communicate to team
4. Document in commit message

- Do not introduce synonyms
- Define terms on first use

If a term exists in the glossary, you MUST use it exactly.

---

### Lists

- Use numbered lists for sequences
- Use bullet lists for collections
- Use parallel grammar

---

### Warnings and Notes

Use the following labels consistently:

- **Note:** Adds context or clarification.
- **Caution:** Indicates potential problems.
- **Warning:** Indicates risk of data loss or failure.

Do not overuse warnings.

---

## 💻 Code & Examples

### Code Blocks

- Always specify the language
- Keep examples minimal and runnable
- Do not include placeholder code unless clearly marked

```bash
npm run docs:dev
```

```

---

### Inline Code

- Use inline code for:

  - Commands
  - Filenames
  - Flags
  - UI labels

---

## 📁 File and Navigation Rules

- One concept per page
- Do not exceed 800 words without justification
- Prefer multiple short pages over one long page
- Sidebar labels must match page titles
- Avoid deep nesting (max 3 levels)

---

## 🔍 SEO and Discoverability

- First paragraph must define the page clearly
- Avoid vague headings
- Use descriptive link text

✅ Good:

> See Configure authentication

❌ Bad:

> Click here

---

## 🚫 Prohibited Content

The following are NOT allowed:

- Marketing language
- Emojis
- Humor
- Sarcasm
- Metaphors
- Personal opinions
- Unverified claims
- Future promises

---

## 🤖 Automated Agents and AI

If an automated agent contributes documentation:

- Output MUST comply with this file
- No conversational language
- No acknowledgements or sign-offs
- No self-references
- No explanations of reasoning

Generated content must be indistinguishable from human-written documentation.

---

## 🧪 Review Checklist

Before merging documentation, verify:

- Style matches Google guidelines
- Grammar and spelling are correct
- Instructions are reproducible
- Links work
- Headings are consistent
- Examples are accurate
- Tone is neutral and professional

If any item fails, request changes.

---

## 📌 Enforcement

Pull requests that do not follow this guide MUST be rejected.

Consistency matters more than author preference.

---

## 📎 References

- Google Developer Documentation Style Guide
  [https://developers.google.com/style](https://developers.google.com/style)
```
