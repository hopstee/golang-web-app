export interface Entity {
	id: string;
	title: string;
	slug: string;
	content: Record<string, unknown>;
	createdAt: string;
	updatedAt: string;
}

export interface ShortEntityData {
	id: string;
	title: string;
}

export const EntityDataTypes = {
	LAYOUT_FIELDS: "layout_fields",
	CONTENT: "content",
} as const

export interface EntityData {
	content: Record<string, unknown>;
	[key: string]: unknown;
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