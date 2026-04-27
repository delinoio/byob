var count = 1;
const unusedValue: any = "demo";

console.log(count == 1);

debugger;

export function invokeLater(callback: () => void) {
	setTimeout(callback, 10);
}
