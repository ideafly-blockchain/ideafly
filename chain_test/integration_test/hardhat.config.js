require("@nomiclabs/hardhat-waffle");
require("hardhat-deploy");
require("hardhat-deploy-ethers");
require("hardhat-contract-sizer");
require("hardhat-gas-reporter");
require("@nomiclabs/hardhat-truffle5");

const accounts = [
    "0xebaa2febee077847f41b9bd23b28ba7318f37d92658ccbe194a2df432a93810f",
    "0xdf504d175ae63abf209bad9dda965310d99559620550e74521a6798a41215f46",
    "0x78e664ccfe34a872dc6f0962a37e3ac77f6980b1ba6ab2d485d2af175a533721",
    "0xce571ffbe8190957564e41bb84348a7f9d47a475639405b90b49c026784d97f5",
    "0xb80f57ce08935019c3fb61f56a9ac9202f3e0a2d53b29999890bb61fc4e0421e",
    "0x95357e3d1c869ff1dee172d3617b8ce74266a3b135236d04fce4363efc76de37",
    "0x52536142e4d0728d16014bfd47628d27b407da138d4f61046c85f109a95547f1",
    "0x80ad06214c1be9cc3cd941d59f4a53ce0d2b7e24fc646fb48f81d8cf5f74a6bc",
    "0xed8da5520a704b8e417c0d3ab2511bb6eebc929d3812b90db7a4ac17132647f2",
    "0x421d0fae37983caee801945267d9ab514fdefe51113389a380b8e1e2de968e4c",
    "0x1aa2a104115ad2e1808000c36349902e557fdfa259e2da9371ce5cc93b1753e2",
    "0x0fbcb1a73abc49a8cce055dede1d3b138e4cebfabc7d58aa90f8e2c662f2f1af",
    "0x2441546bd4f709be6cb83f2aa177b0bcc00130e900a9c21b0b8837aa3988fd1f",
    "0x9a02ae5cb5967cc8a27d32f8059988ff1b0a554891608147dffc9bd8e6ee0633",
    "0xb07907317e4a5c2731ebdf7ba781eeb28db1630096d7998755612a3b5a8d118b",
    "0x5e9561af4f2963911d4c04c0fe830666f57b0d87f9bd24ffc4f65aad2a2c2de1",
];
/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
    solidity: {
        compilers: [
            {
                version: "0.8.20",
                settings: {
                    evmVersion: "berlin",
                },
            },
        ],
    },
    defaultNetwork: "local",
    networks: {
        local: {
            chainId: 6660001,
            url: "http://127.0.0.1:8545",
            accounts,
            gasPrice: 2000000000,
        },
    },
    contractSizer: {
        alphaSort: true,
        runOnCompile: true,
        disambiguatePaths: false,
    },
    spdxLicenseIdentifier: {
        overwrite: true,
        runOnCompile: true,
    },
};
