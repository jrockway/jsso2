describe("Enrolls, logs in, visits a protected web page", () => {
    it("Can't yet visit a protected webpage", () => {
        cy.clearCookie("jsso-session-id");
        cy.visit("/protected");
        cy.url().should("include", "/#/login");
    });
    it("Enrolls a security key", () => {
        cy.clearCookie("jsso-session-id");
        cy.task("addKey").then(() => {
            cy.task("enrollmentLink").then((enrollmentLink: string) => {
                cy.visit(enrollmentLink);
                cy.get(".error").should("not.exist");
                cy.contains("Welcome, the-tests");
                cy.get("#enroll")
                    .click()
                    .then(() => {
                        cy.contains("Successfully added").should("exist");
                    });
            });
        });
    });
    it("Can login", () => {
        cy.clearCookie("jsso-session-id");
        cy.visit("/#/login");
        cy.get("#username")
            .type("the-tests")
            .then(() => {
                cy.get("#login")
                    .click()
                    .then(() => {
                        cy.url().should("equal", "http://localhost:4000/");
                        cy.should("not.contain", "Proceed to the login page");
                        cy.contains("Welcome, the-tests").should("exist");
                    });
            });
    });
    it("Can visit a protected page", () => {
        cy.visit("/protected");
        cy.contains("Logged in as the-tests").should("exist");
    });
    it("Can logout", () => {
        cy.visit("/logout");
        cy.url().should("contain", "/#/login");
    });
    it("Can no longer see the protected page", () => {
        cy.visit("/protected");
        cy.url().should("include", "/#/login");
    });
});
