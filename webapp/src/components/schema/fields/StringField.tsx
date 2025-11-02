import { Input } from "@/components/ui/input";
import type { Field } from "@/types/pages";

interface StringField {
    field: Field;
    value: string;
    onChange: (name: string, value: string) => void;
}

export default function StringField(props: StringField) {
    return (
        <Input
            value={props.value}
            onChange={(e) => props.onChange(props.field.name, e.target.value)}
            placeholder={props.field.label}
        />
    );
}