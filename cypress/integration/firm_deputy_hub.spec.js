describe("Firm Deputy Hub", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/firm/");
  });

    it('finds the content "Hello world!"', () => {
        cy.contains('Hello world!')
    })
});
