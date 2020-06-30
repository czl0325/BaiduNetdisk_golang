function decodeQuery() {
    const search = decodeURI(document.location.search);
    return search.replace(/(^\?)/, '').split('&').reduce(function (result, item) {
        let values = item.split('=');
        result[values[0]] = values[1];
        return result;
    }, {});
}

function getCookie(name) {
    const r = document.cookie.match("\\b" + name + "=([^;]*)\\b");
    return r ? r[1] : undefined;
}