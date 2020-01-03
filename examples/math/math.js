var nf = nf || {};
nf.math = {};
nf.util = {};

nf.list = function(type, cnt) {
    //TODO
}

nf.it = function() {}
nf.what = function() {}

nf.math.expression = function(s) {
    return s.split("").join('*');
}

nf.math.to_number = function(s) {
    return Number(s);
}

nf.math.decimal = function(s) {
    s = s.toString();
    var n = Number(s);
    return n / Math.pow(10, s.length);
}

nf.math.sum = function(x, y) {
    return x + y;
}

nf.math.sub = function(x, y) {
    return x - y;
}

nf.math.mul = function(x, y) {
    return x * y;
}

nf.math.div = function(x, y) {
    return x / y;
}

nf.math.neg = function(x) {
    return -x;
}

nf.math.pow = function(x, y) {
    return Math.pow(x, y);
}

nf.util.concat = function(x, y) {
    return x.toString() + y.toString();
}
