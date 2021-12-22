describe("Manage Pii Details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/firm/1/manage-pii-details");
        });

        it("will show a validation error if form submitted with empty fields", () => {
            cy.setCookie("fail-route", "pii-details");
            cy.get("[data-cy=submit-manage-pii-details-form-btn]").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );

            cy.get(".govuk-error-summary").should(
                "contain",
                "The pii expiry date is required and can't be empty"
            ); 
            
            cy.get(".govuk-error-summary").should(
                "contain",
                "The pii received date is required and can't be empty"
            );

            cy.get(".govuk-error-summary").should(
                "contain",
                "The pii amount is required and can't be empty"
            );
        });
    });
})