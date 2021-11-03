module.exports = {
    ci: {
      collect: {
        url: ['http://firm-deputy-hub:8888/supervision/deputies/firm/'],
        settings: {
          extraHeaders: JSON.stringify({Cookie: 'XSRF-TOKEN=abcde; Other=other'}),
          chromeFlags: "--disable-gpu --no-sandbox",
        },
      },
      assert: {
        assertions: {
          "categories:performance": ["warn", {"minScore": 0.85}],
          "categories:accessibility": ["error", {"minScore": 0.85}],
          "categories:best-practices": ["warn", {"minScore": 0.9}],
          "categories:seo": ["warn", {"minScore": 0.7}],
        },
      },
      upload: {
        target: 'temporary-public-storage',
      },
    },
  };
  