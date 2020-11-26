/// <reference types="cypress" />
// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)
const CRI = require("chrome-remote-interface");

let criClient = null;
let criPort = 0;

/**
 * @type {Cypress.PluginConfig}
 */
module.exports = (on, config) => {
    on("before:browser:launch", (browser, args) => {
        args = require("cypress-log-to-output").browserLaunchHandler(browser, args);
        criPort = ensureRdpPort(args.args);
        console.log("criPort is", criPort);
    });
    on("task", {
        enrollmentLink() {
            return process.env.ENROLLMENT_LINK;
        },
        async addKey() {
            criClient = criClient || (await CRI({ port: criPort }));
            await criClient.send("WebAuthn.enable", {});
            return await criClient.send("WebAuthn.addVirtualAuthenticator", {
                options: {
                    protocol: "ctap2",
                    transport: "usb",
                    hasResidentKey: false,
                    hasUserVerification: true,
                    isUserConsenting: true,
                    isUserVerified: true,
                },
            });
        },
        async cri(args) {
            criClient = criClient || (await CRI({ port: criPort }));
            return criClient.send(args.query, args.opts);
        },
    });
};

function ensureRdpPort(args) {
    const existing = args.find((arg) => arg.slice(0, 23) === "--remote-debugging-port");

    if (existing) {
        return Number(existing.split("=")[1]);
    }

    const port = 40000 + Math.round(Math.random() * 25000);
    args.push(`--remote-debugging-port=${port}`);
    return port;
}
