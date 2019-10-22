import {MessageTypes} from "./messages";

const recentSearches = [
    "anthony kiedis",
    "kurt vonnegut",
    "dune frank herbert long title test"
];


const servers = {
    "type": MessageTypes.SERVERS,
    "servers": [
        "LawdyServer",
        "Oatmeal",
        "Astra"
    ]
};

const fakeItems = {
    "type": 2,
    "books": [{
        "server": "LawdyServer",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (html)",
        "format": "html",
        "size": "151.8KB",
        "full": "!LawdyServer Kurt Vonnegut - The Sirens of Titan (v5.0) (html).rar "
    }, {
        "server": "LawdyServer",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (mobi)",
        "format": "mobi",
        "size": "269.8KB",
        "full": "!LawdyServer Kurt Vonnegut - The Sirens of Titan (v5.0) (mobi).zip "
    }, {
        "server": "LawdyServer",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!LawdyServer Kurt Vonnegut - The Sirens of Titan (v5.0).epub "
    }, {
        "server": "LawdyServer",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan [pdf]",
        "format": "pdf",
        "size": "862.7KB",
        "full": "!LawdyServer Kurt Vonnegut - The Sirens of Titan [pdf].pdf "
    }, {
        "server": "LawdyServer",
        "author": "Kurt Vonnegut Jr",
        "title": "The Sirens of Titan (v5.0)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!LawdyServer Kurt Vonnegut Jr - The Sirens of Titan (v5.0).epub "
    }, {
        "server": "LawdyServer",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan",
        "format": "pdf",
        "size": "862.7KB",
        "full": "!LawdyServer Kurt Vonnegut - The Sirens of Titan.pdf "
    }, {
        "server": "shytot",
        "author": "The sirens of Titan",
        "title": "Kurt Vonnegut",
        "format": "epub",
        "size": "339.9KB",
        "full": "!shytot The sirens of Titan - Kurt Vonnegut.epub "
    }, {
        "server": "shytot",
        "author": "The Sirens of Titan",
        "title": "Kurt Vonnegut",
        "format": "mobi",
        "size": "365.3KB",
        "full": "!shytot The Sirens of Titan - Kurt Vonnegut.mobi "
    }, {
        "server": "shytot",
        "author": "The Sirens of Titan",
        "title": "Kurt Vonnegut",
        "format": "epub",
        "size": "250.1KB",
        "full": "!shytot The Sirens of Titan - Kurt Vonnegut.epub "
    }, {
        "server": "shytot",
        "author": "The Sirens of Titan",
        "title": "Kurt Vonnegut",
        "format": "mobi",
        "size": "365.3KB",
        "full": "!shytot The Sirens of Titan - Kurt Vonnegut.mobi "
    }, {
        "server": "Wench",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!Wench Kurt Vonnegut - The Sirens of Titan (v5.0) (epub).rar "
    }, {
        "server": "Wench",
        "author": "Kurt Vonnegut",
        "title": "The Sirens Of Titan",
        "format": "mobi",
        "size": "378.6KB",
        "full": "!Wench Kurt Vonnegut - The Sirens Of Titan.mobi "
    }, {
        "server": "DukeLupus",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (html)",
        "format": "html",
        "size": "151.84KB",
        "full": "!DukeLupus Kurt Vonnegut - The Sirens of Titan (v5.0) (html).rar"
    }, {
        "server": "DukeLupus",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (mobi)",
        "format": "mobi",
        "size": "269.78KB",
        "full": "!DukeLupus Kurt Vonnegut - The Sirens of Titan (v5.0) (mobi).zip"
    }, {
        "server": "DukeLupus",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0)",
        "format": "epub",
        "size": "1.83MB",
        "full": "!DukeLupus Kurt Vonnegut - The Sirens of Titan (v5.0).epub"
    }, {
        "server": "DukeLupus",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan [pdf]",
        "format": "pdf",
        "size": "862.67KB",
        "full": "!DukeLupus Kurt Vonnegut - The Sirens of Titan [pdf].pdf"
    }, {
        "server": "DukeLupus",
        "author": "Kurt Vonnegut Jr",
        "title": "The Sirens of Titan (v5.0)",
        "format": "epub",
        "size": "1.83MB",
        "full": "!DukeLupus Kurt Vonnegut Jr - The Sirens of Titan (v5.0).epub"
    }, {
        "server": "DukeLupus",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan",
        "format": "pdf",
        "size": "862.67KB",
        "full": "!DukeLupus Kurt Vonnegut - The Sirens of Titan.pdf"
    }, {
        "server": "dragnbreaker",
        "author": "Vonnegut, Kurt",
        "title": "Novel 02 - The Sirens of Titan (v5.0)",
        "format": "epub",
        "size": "282.5KB",
        "full": "!dragnbreaker Vonnegut, Kurt - Novel 02 - The Sirens of Titan (v5.0).epub "
    }, {
        "server": "dragnbreaker",
        "author": "Hugo 1960 Novel Nominee",
        "title": "Kurt Vonnegut - The Sirens of Titan",
        "format": "rtf",
        "size": "554.9KB",
        "full": "!dragnbreaker Hugo 1960 Novel Nominee - Kurt Vonnegut - The Sirens of Titan.rtf "
    }, {
        "server": "dragnbreaker",
        "author": "Masterworks",
        "title": "SF Masterworks 018 - The Sirens of Titan - Kurt Vonnegut",
        "format": "epub",
        "size": "250.1KB",
        "full": "!dragnbreaker Masterworks - SF Masterworks 018 - The Sirens of Titan - Kurt Vonnegut.epub "
    }, {
        "server": "Astra",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (epub)",
        "format": "epub",
        "size": "250.1KB",
        "full": "!Astra Kurt Vonnegut - The Sirens of Titan (epub).epub "
    }, {
        "server": "Astra",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!Astra Kurt Vonnegut - The Sirens of Titan (v5.0) (epub).epub "
    }, {
        "server": "Astra",
        "author": "Kurt Vonnegut Jr",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!Astra Kurt Vonnegut Jr - The Sirens of Titan (v5.0) (epub).epub "
    }, {
        "server": "bsk",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (epub)",
        "format": "epub",
        "size": "251.3KB",
        "full": "!bsk Kurt Vonnegut - The Sirens of Titan (epub).rar "
    }, {
        "server": "bsk",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (html)",
        "format": "html",
        "size": "151.8KB",
        "full": "!bsk Kurt Vonnegut - The Sirens of Titan (v5.0) (html).rar "
    }, {
        "server": "bsk",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (mobi)",
        "format": "mobi",
        "size": "270.1KB",
        "full": "!bsk Kurt Vonnegut - The Sirens of Titan (v5.0) (mobi).rar "
    }, {
        "server": "bsk",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!bsk Kurt Vonnegut - The Sirens of Titan (v5.0) (epub).rar "
    }, {
        "server": "bsk",
        "author": "Kurt Vonnegut Jr",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!bsk Kurt Vonnegut Jr - The Sirens of Titan (v5.0) (epub).rar "
    }, {
        "server": "bsk",
        "author": "057",
        "title": "Vonnegut, Kurt Jr. - The Sirens Of Titan (V5.0) [doc, txt]",
        "format": "rar",
        "size": "1.9MB",
        "full": "!bsk 057 - Vonnegut, Kurt Jr. - The Sirens Of Titan (V5.0) [doc, txt].rar "
    }, {
        "server": "Oatmeal",
        "author": "057",
        "title": "Vonnegut, Kurt Jr. - The Sirens Of Titan (V5.0) [doc, txt]",
        "format": "rar",
        "size": "1.9MB",
        "full": "!Oatmeal 057 - Vonnegut, Kurt Jr. - The Sirens Of Titan (V5.0) [doc, txt].rar "
    }, {
        "server": "Oatmeal",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (html)",
        "format": "html",
        "size": "133.9KB",
        "full": "!Oatmeal Kurt Vonnegut - The Sirens of Titan (v5.0) (html).rar "
    }, {
        "server": "Oatmeal",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (mobi)",
        "format": "mobi",
        "size": "266.5KB",
        "full": "!Oatmeal Kurt Vonnegut - The Sirens of Titan (v5.0) (mobi).rar "
    }, {
        "server": "Oatmeal",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (epub)",
        "format": "epub",
        "size": "251.3KB",
        "full": "!Oatmeal Kurt Vonnegut - The Sirens of Titan (epub).rar "
    }, {
        "server": "Oatmeal",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!Oatmeal Kurt Vonnegut - The Sirens of Titan (v5.0) (epub).rar "
    }, {
        "server": "Oatmeal",
        "author": "Kurt Vonnegut Jr",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!Oatmeal Kurt Vonnegut Jr - The Sirens of Titan (v5.0) (epub).rar "
    }, {
        "server": "phoomphy",
        "author": "Vonnegut, Kurt",
        "title": "Sirens of Titan (1959)",
        "format": "epub",
        "size": "282.49 KiB",
        "full": "!phoomphy Vonnegut, Kurt - Sirens of Titan (1959).epub    "
    }, {
        "server": "Pondering42",
        "author": "057",
        "title": "Vonnegut, Kurt Jr. - The Sirens Of Titan (V5.0) [doc, txt]",
        "format": "rar",
        "size": "1.87MB",
        "full": "!Pondering42 057 - Vonnegut, Kurt Jr. - The Sirens Of Titan (V5.0) [doc, txt].rar"
    }, {
        "server": "Pondering42",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (html)",
        "format": "html",
        "size": "151.84KB",
        "full": "!Pondering42 Kurt Vonnegut - The Sirens of Titan (v5.0) (html).rar"
    }, {
        "server": "Pondering42",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (mobi)",
        "format": "mobi",
        "size": "270.12KB",
        "full": "!Pondering42 Kurt Vonnegut - The Sirens of Titan (v5.0) (mobi).rar"
    }, {
        "server": "Pondering42",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (epub)",
        "format": "epub",
        "size": "251.27KB",
        "full": "!Pondering42 Kurt Vonnegut - The Sirens of Titan (epub).rar"
    }, {
        "server": "Pondering42",
        "author": "Kurt Vonnegut",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!Pondering42 Kurt Vonnegut - The Sirens of Titan (v5.0) (epub).rar"
    }, {
        "server": "Pondering42",
        "author": "Kurt Vonnegut Jr",
        "title": "The Sirens of Titan (v5.0) (epub)",
        "format": "epub",
        "size": "1.8MB",
        "full": "!Pondering42 Kurt Vonnegut Jr - The Sirens of Titan (v5.0) (epub).rar"
    }, {
        "server": "QuietSilence",
        "author": "Vonnegut, Kurt",
        "title": "The Sirens Of Titan [rtf]",
        "format": "rtf",
        "size": "247.9KB",
        "full": "!QuietSilence Vonnegut, Kurt - The Sirens Of Titan [rtf].rar "
    }]
};

export {
    recentSearches,
    fakeItems,
    servers
}