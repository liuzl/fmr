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

var overrideMethod = function() {
	switch (arguments.length) {
		case 0:
			console.log("no arguments");
			break;
		case 1:
			console.log("1");
			console.log(arguments[0]);
			break;
		default:
			for (i = 0; i < arguments.length; i++) {
				console.log(typeof(arguments[i]));
				console.log(arguments[i]);
			}
			break;
	}
}

overrideMethod();
overrideMethod("haha");
overrideMethod("haha", "hehe", 1, 1, 1, [2, 3, 4]);
