export type BuildPrimitive = string | number | boolean | null;

export type BuildValue =
  | BuildPrimitive
  | readonly BuildValue[]
  | { readonly [key: string]: BuildValue };

export type BuildInputKind = "source" | "config" | "asset";

export interface BuildInput {
  readonly path: string;
  readonly kind?: BuildInputKind;
}

export type BuildOutputKind = "artifact" | "metadata" | "diagnostic";

export interface BuildOutput {
  readonly path: string;
  readonly kind?: BuildOutputKind;
}

export interface BuildTaskContext<TInput extends BuildValue = BuildValue> {
  readonly cwd: string;
  readonly input: TInput;
  readonly env: Readonly<Record<string, string | undefined>>;
  readonly readFile?: (path: string) => Promise<string>;
  readonly writeFile?: (path: string, contents: string) => Promise<void>;
  readonly log?: (message: string, fields?: Readonly<Record<string, BuildValue>>) => void;
}

export type BuildTaskAction<
  TInput extends BuildValue = BuildValue,
  TOutput extends BuildValue = BuildValue,
> = (context: BuildTaskContext<TInput>) => TOutput | Promise<TOutput>;

export interface BuildTask<
  TInput extends BuildValue = BuildValue,
  TOutput extends BuildValue = BuildValue,
> {
  readonly description?: string;
  readonly dependsOn?: readonly string[];
  readonly inputs?: readonly BuildInput[];
  readonly outputs?: readonly BuildOutput[];
  readonly run?: BuildTaskAction<TInput, TOutput>;
}

export type BuildTaskMap = Readonly<Record<string, BuildTask>>;

export interface BuildToolDefinition<TTasks extends BuildTaskMap = BuildTaskMap> {
  readonly name: string;
  readonly version?: string;
  readonly tasks: TTasks;
}

export function defineBuildTool<const TTasks extends BuildTaskMap>(
  definition: BuildToolDefinition<TTasks>,
): BuildToolDefinition<TTasks> {
  return definition;
}

