"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Home;
var Link_1 = require("@docusaurus/Link");
var useDocusaurusContext_1 = require("@docusaurus/useDocusaurusContext");
var Heading_1 = require("@theme/Heading");
var Layout_1 = require("@theme/Layout");
var clsx_1 = require("clsx");
var index_module_css_1 = require("./index.module.css");
function HomepageHeader() {
    var siteConfig = (0, useDocusaurusContext_1.default)().siteConfig;
    return (<header className={(0, clsx_1.default)('hero hero--primary', index_module_css_1.default.heroBanner)}>
      <div className="container">
        <img src="img/logo.svg" width="200"/>
        <Heading_1.default as="h1" className="hero__title">
          {siteConfig.title}
        </Heading_1.default>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={index_module_css_1.default.buttons}>
          <Link_1.default className="button button--secondary button--lg" to="/docs/intro">
            Read the handbook
          </Link_1.default>
        </div>
      </div>
    </header>);
}
function Home() {
    var siteConfig = (0, useDocusaurusContext_1.default)().siteConfig;
    return (<Layout_1.default title={"Welcome to ".concat(siteConfig.title)} description="Operational guidance, product practices, and cultural guardrails for AVM.">
      <HomepageHeader />
    </Layout_1.default>);
}
