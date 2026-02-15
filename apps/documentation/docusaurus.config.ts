import type * as Preset from "@docusaurus/preset-classic";
import type { Config } from "@docusaurus/types";
import { themes as prismThemes } from "prism-react-renderer";

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: "CDNS Documentation",
  tagline: "A trusted, Linux-first DNS management CLI tool",
  favicon: "img/favicon.ico",

  // Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // Improve compatibility with the upcoming Docusaurus v4
  },

  // Set the production url of your site here
  url: "https://junevm.gitlab.io",
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/cdns/",

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "junevm", // Used for metadata and default edit links.
  projectName: "cdns",

  onBrokenLinks: "throw",

  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      {
        pages: {
          exclude: ["**/*.js"],
        },
        docs: {
          sidebarPath: "./sidebars.ts",
        },
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themes: ["@docusaurus/theme-mermaid"],

  markdown: {
    mermaid: true,
  },

  themeConfig: {
    // Replace with your project's social card
    image: "img/logo.svg",
    colorMode: {
      defaultMode: "dark",
      respectPrefersColorScheme: false,
    },
    navbar: {
      title: "CDNS",
      logo: {
        alt: "CDNS Logo",
        src: "img/logo.svg",
      },
      items: [
        {
          type: "docSidebar",
          sidebarId: "docsSidebar",
          position: "left",
          label: "Documentation",
        },
        {
          href: "https://gitlab.com/junevm/cdns",
          label: "Repository",
          position: "right",
        },

        {
          type: "search",
          position: "right",
        },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Documentation",
          items: [
            {
              label: "Introduction",
              to: "/docs/intro",
            },
          ],
        },
        {
          title: "Community",
          items: [
            {
              label: "GitLab",
              href: "https://gitlab.com/junevm/cdns",
            },
          ],
        },
      ],
    },
    mermaid: {
      theme: { light: "base", dark: "base" },
      options: {
        fontFamily: "var(--ifm-font-family-base)",
        fontSize: 14,
        themeVariables: {
          // Twitter-inspired color scheme
          primaryColor: "#1DA1F2",
          primaryTextColor: "#FFFFFF",
          primaryBorderColor: "#1A91DA",
          lineColor: "#1DA1F2",
          sectionBkgColor: "#192734",
          altSectionBkgColor: "#22303C",
          gridColor: "#2F3336",
          tertiaryColor: "#3D5466",
          background: "#15202B",
          mainBkg: "#192734",
          secondBkg: "#22303C",
          border1: "#2F3336",
          border2: "#3D5466",
          textColor: "#FFFFFF",
          tertiaryTextColor: "#8B98A5",
          noteBkgColor: "#22303C",
          noteTextColor: "#FFFFFF",
          noteBorderColor: "#2F3336",
          labelTextColor: "#FFFFFF",
          errorBkgColor: "#F4212E",
          errorTextColor: "#FFFFFF",
        },
      },
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
