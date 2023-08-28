function readFile(filename) {
    const data = require('fs').readFileSync(filename,'utf-8')
    return data.split('\n')
}

function linearToRGB(linear) {
    return (255.999 * Math.min(1, Math.pow(linear,1.0/2.2))) | 0
	return (255.999 * Math.sqrt(Math.min(1, linear))) | 0
}


const data = [
    readFile('/Users/alessandro/Downloads/i21_1.ppm'),
    readFile('/Users/alessandro/Downloads/i21_2.ppm'),
    readFile('/Users/alessandro/Downloads/i21_3.ppm')
]

console.log('P3')
console.log('400 400')
console.log('255')

for( let i=3; i<data[0].length; i++ ) {
    if( data[0][i].length < 4 ) {
        console.log('')
    }
    else {
        let r = 0, b = 0, g = 0

        for(let j=0; j<2; j++) {
            rgb = data[j][i].split(' ')
            r += parseFloat(rgb[0])
            g += parseFloat(rgb[1])
            b += parseFloat(rgb[2])
        }

        const f = 3

        console.log(linearToRGB(r / f), linearToRGB(g / f), linearToRGB(b / f))
    }
}