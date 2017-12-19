var nf = nf || {};
nf.math = {};

nf.math.sum = function(x, y) {
    console.log(typeof(x), typeof(y));
    console.log(x, y);
    return x + y;
};
nf.math.minus = function(x, y) {
    return x - y;
}
nf.math.times = function(x, y) {
    return x * y;
}
nf.math.neg = function(x) {
    return -x;
}

//console.log(nf.math.sum(1,1));
