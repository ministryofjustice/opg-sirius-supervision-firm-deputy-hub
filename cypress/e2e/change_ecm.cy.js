describe("Change Ecm", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        it("has a link from the main firms page", () => {
            cy.visit("/supervision/deputies/firm/1");
            cy.get(".moj-button-menu__wrapper > .govuk-button")
                .should("contain", "Change ECM")
                .click();
            cy.url().should("include", "firm/1/change-ecm");
            cy.get(".govuk-heading-l").should(
                "contain",
                "Change Executive Case Manager"
            );
        });

        it("has a cancel link which returns to the main firms page", () => {
            cy.visit("/supervision/deputies/firm/1/change-ecm");
            cy.get("#change-ecm-cancel").should("contain", "Cancel").click();
            cy.url().should("not.include", "firm/1/change-ecm");
            cy.get(".govuk-heading-l").should("contain", "Firm details");
        });
    });

    describe("Change Ecm Form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/firm/1/change-ecm");
        });

        it("should autopopulate the existing firm data", () => {
            cy.get(".govuk-body").should("contain", "LayTeam1 User6");
        });

        it("should have team data in autocomplete", () => {
            cy.get("#select-ecm").type("L");
            cy.get("#select-ecm__listbox").find("li").should("have.length", 3);
            cy.get("#select-ecm").type("uke");
            cy.get("#select-ecm__listbox").find("li").should("have.length", 1);
        });

        it("allows me to successfully submit the form", () => {
            cy.setCookie("success-route", "changeECM");
            cy.get("#select-ecm").type("Han Solo");
            cy.get(".govuk-button").should("contain", "Change ECM").click();
            cy.url().should("not.include", "firm/1/change-ecm");
            cy.get(".moj-banner").should(
                "contain",
                "ECM changed to LayTeam1 User6"
            );
            cy.get(".govuk-heading-l").should("contain", "Firm details");
        });

        it("will show a validation error if form submitted with empty fields", () => {
            cy.get(".govuk-button").should("contain", "Change ECM").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-list > li").should(
                "contain",
                "Select an executive case manager"
            );
        });
    });
});
