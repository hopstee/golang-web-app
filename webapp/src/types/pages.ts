export interface Page {
	id: string;
	title: string;
	slug: string;
	content: PageDataData;
	createdAt: string;
	updatedAt: string;
}

export const PageDataTypes = {
	LAYOUT_FIELDS: "layout_fields",
	CONTENT: "content",
} as const

export interface PageDataData {
	layout_fields: Record<string, unknown>
	content: Record<string, unknown>
}

export interface PageData {
	data: PageDataData;
	schema: {
		id: string;
		title: string;
		type: string;
		layout: string;
		layout_fields: Field[];
		content: Field[];
	}
}

export interface Field {
	id: string;
	name: string;
	type: string;
	label: string;
	schema: Schema;
}

export interface Schema {
	fields: Field[];
}

export const PageSchemaType = {
	LAYOUT: "layout",
	BLOCK: "block",
	PAGE: "page",
	MODULE: "module",
	SHARED: "shared",
} as const

export interface PageSchema {
	id: string;
	title: string;
	type: typeof PageSchemaType;
	layout: string;
	parent: string;
	blocks: string[];
	seo: Field[];
	content: Field[];
	children: PageSchema[];
}