import type { Field } from "@/types/entities";
import { Label } from "@/components/ui/label";
import { useEffect, useState } from "react";
import FileUploader from "@/components/schema/fields/FileUploader";
import StringField from "@/components/schema/fields/StringField";
import TextField from "@/components/schema/fields/TextField";
import BoolField from "@/components/schema/fields/BoolField";
import ListField from "@/components/schema/fields/ListField";
import { getDefaultValue } from "@/lib/utils";
import SelectField from "@/components/schema/fields/SelectField";
import NumberField from "@/components/schema/fields/NumberField";

interface FieldsScaffoldProps {
    fields: Field[];
    data: Record<string, unknown>
    onChange?: (updated: Record<string, unknown>) => void
}

export default function FieldsScaffold(props: FieldsScaffoldProps) {
    const {
        fields,
        data,
        onChange
    } = props;

    const [values, setValues] = useState<Record<string, unknown>>(() =>
        fields.reduce((acc, f) => ({ ...acc, [f.name]: data[f.name] ?? getDefaultValue(f.type) }), {})
    )

    useEffect(() => {
        setValues(
            fields.reduce(
                (acc, f) => ({ ...acc, [f.name]: data[f.name] ?? getDefaultValue(f.type) }),
                {}
            )
        )
    }, [data, fields])

    const handleChange = (name: string, value: unknown) => {
        const updated = { ...values, [name]: value }
        setValues(updated);
        onChange?.(updated);
    }

    const handleDelete = (name: string, index: number) => {
        const list = [...((values[name] as unknown[]) || [])]
        list.splice(index, 1);
        handleChange(name, list);
    }

    const renderField = (field: Field) => {
        const value = values[field.name] ?? ""

        switch (field.type) {
            case "string":
                return <StringField field={field} value={value as string} onChange={handleChange} />;
            case "number":
                return <NumberField field={field} value={Number(value)} onChange={handleChange} />;
            case "text":
                return <TextField field={field} value={value as string} onChange={handleChange} />;
            case "select":
                return <SelectField field={field} value={value as string} onChange={handleChange} />;
            case "bool":
            case "boolean":
                return <BoolField field={field} value={String(value)} onChange={handleChange} />;
            case "img":
                return <FileUploader field={field} value={value as string} onChange={handleChange} />;
            case "list[string]":
                return <ListField
                    type="list[string]"
                    field={field}
                    value={value as unknown[]}
                    handleChange={handleChange}
                    handleDelete={handleDelete}
                />;
            case "list[text]":
                return <ListField
                    type="list[text]"
                    field={field}
                    value={value as unknown[]}
                    handleChange={handleChange}
                    handleDelete={handleDelete}
                />;

            case "list[object]":
                return (
                    <ListField
                        type="list[object]"
                        field={field}
                        value={value as unknown[]}
                        handleChange={handleChange}
                        handleDelete={handleDelete}
                    />
                );
        }
    }

    if (fields.length === 0) {
        return <div>Нет доступных полей для редактирования</div>
    }

    console.log(fields)
    console.log(values)

    return (
        <div className="space-y-6">
            {fields.map(field => {
                if (field.depends && !field.depends.values.includes(values[field.depends.field] as string)) {
                    return;
                }
                
                return (
                    <div className="space-y-3" key={field.id}>
                        <Label key={field.id}>{field.label}</Label>
                        {renderField(field)}
                    </div>
                )
            })}
        </div>
    )
}
