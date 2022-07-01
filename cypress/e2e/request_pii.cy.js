describe("Request PII Details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/firm/1/request-pii-details");
        });

        it("can be submitted", () => {
            cy.setCookie("success-route", "requestPii");
            cy.get("#f-pii-requested").type("2013-12-12")
            cy.get("[data-cy=submit-request-pii-details-form-btn]").click();
        });

        it("will show a validation error if form submitted with empty fields", () => {
            cy.setCookie("fail-route", "request-pii-details");
            cy.get("[data-cy=submit-request-pii-details-form-btn]").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );

            cy.get(".govuk-error-summary").should(
                "contain",
                "The PII requested date is required and can't be empty"
            );
        });

        it("will autofill the form fields", () => {
            cy.get("#f-pii-requested").should(
                "have.value",
                "2000-11-15"
            );
        });
    });
})