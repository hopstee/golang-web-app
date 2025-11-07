export interface Entity {
	id: string;
	title: string;
	slug: string;
	content: EntityDataData;
	createdAt: string;
	updatedAt: string;
}

export const EntityDataTypes = {
	LAYOUT_FIELDS: "layout_fields",
	CONTENT: "content",
} as const

export interface EntityDataData {
	layout_fields: Record<string, unknown>
	content: Record<string, unknown>
}

export interface EntityData {
	data: EntityDataData;
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

export const EntitySchemaType = {
	LAYOUT: "layout",
	BLOCK: "block",
	PAGE: "page",
	MODULE: "module",
	SHARED: "shared",
} as const

export interface EntitySchema {
	id: string;
	title: string;
	type: typeof EntitySchemaType;
	layout: string;
	parent: string;
	blocks: string[];
	seo: Field[];
	content: Field[];
	children: EntitySchema[];
}