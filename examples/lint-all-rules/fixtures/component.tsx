declare namespace JSX {
	interface IntrinsicElements {
		button: { type?: string; children?: unknown };
		li: { children?: unknown };
		ul: { children?: unknown };
	}
}

export function List() {
	return (
		<ul>
			{["alpha", "beta"].map((item) => (
				<li>{item}</li>
			))}
			<button>Save</button>
		</ul>
	);
}
