export interface Page {
    id: string;
    title: string;
    slug: string;
    content: Record<string, unknown>;
    createdAt: string;
    updatedAt: string;
}

export interface PageData {
    [key: string]: unknown
}

export interface Field {
	name: string
	type: string
	label: string
	schema: Schema
}

export interface Schema {
	fields: Field[]
}

export const PageSchemaType = {
    LAYOUT: "layout",
    BLOCK: "block",
    PAGE: "page",
    MODULE: "module",
} as const

export interface PageSchema {
	id: string
	title: string
	type: typeof PageSchemaType
	layout: string
	parent: string
	blocks: string[]
	seo: Field[]
	content: Field[]
	children: PageSchema[]
}