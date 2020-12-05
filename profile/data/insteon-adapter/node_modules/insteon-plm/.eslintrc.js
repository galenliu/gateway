module.exports = {
    "env": {
        "node": true,
        "es6": true,
        "mocha": true
    },
    "extends": "eslint:recommended",
    "parserOptions": {
        "ecmaVersion": 2018,
        "sourceType": "module"
    },
    "rules": {
        "no-constant-condition": ["error", { "checkLoops": false }],
        "eol-last": ["error", "always"],
        "class-methods-use-this": ["error"],
        "indent": ["error", 2],
        "brace-style": ["error"],
        "key-spacing": ["error"],
        "keyword-spacing": ["error"],
        "object-curly-spacing": ["error", "always"],
        "object-curly-newline": ["error"]
    }
};
